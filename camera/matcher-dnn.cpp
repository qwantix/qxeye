
#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <fstream>
#include <stdio.h>
#include <opencv2/opencv.hpp>
#include <opencv2/dnn.hpp>

#include "camera.hpp"
#include "matcher.hpp"
#include "matcher-dnn.hpp"
#include "opencv-util.hpp"
#include "util.hpp"

using namespace std;
using namespace cv;
using namespace cv::dnn;

void MatcherDnn::detectMatchesOnRoi(Camera *cam, Rect roi, vector<Match> *matches)
{
  Mat frame = this->frame(roi);
  const bool inverseRgb = frame.channels() == 4;
  Mat inputBlob = blobFromImage(frame,
                                this->inScaleFactor, this->inputBlobSize, this->meanVal, inverseRgb);
  this->net.setInput(inputBlob, this->netInput);
  Mat prob = this->net.forward(this->netOutput);
  Mat detectionMat(prob.size[2], prob.size[3], CV_32F, prob.ptr<float>());
  for (int i = 0; i < detectionMat.rows; i++)
  {
    float confidence = detectionMat.at<float>(i, 2);
    if (confidence > this->minConfidence)
    {
      uint objectClass = (uint)(detectionMat.at<float>(i, 1));
      int xLeftBottom = static_cast<int>(detectionMat.at<float>(i, 3) * frame.cols);
      int yLeftBottom = static_cast<int>(detectionMat.at<float>(i, 4) * frame.rows);
      int xRightTop = static_cast<int>(detectionMat.at<float>(i, 5) * frame.cols);
      int yRightTop = static_cast<int>(detectionMat.at<float>(i, 6) * frame.rows);
      string name = this->classes[to_string(objectClass)];
      Match feature;
      feature.name = name;
      feature.confidence = confidence;
      feature.roi = Rect((int)xLeftBottom + roi.x, (int)yLeftBottom + roi.y,
                         (int)(xRightTop - xLeftBottom),
                         (int)(yRightTop - yLeftBottom));
      matches->push_back(feature);
    }
  }
}

void MatcherDnn::init()
{
  // Init settings
  this->config->get("modelTxt", this->modelTxt);
  this->config->get("modelBin", this->modelBin);
  this->config->get("netInput", this->netInput);
  this->config->get("netOutput", this->netOutput);
  this->config->get("inScaleFactor", this->inScaleFactor);
  this->config->get("meanVal", this->meanVal);
  this->config->get("roiPadding", this->roiPadding);
  this->config->get("minConfidence", this->minConfidence);

  int inputSize = this->inputBlobSize.width;
  this->config->get("inputSize", inputSize);
  this->inputBlobSize.width = this->inputBlobSize.height = inputSize;

  if (this->modelBin.find(".caffemodel") != string::npos)
  {
    this->net = readNetFromCaffe(this->modelTxt, this->modelBin);
  }
  // TODO load YOLO darknet
  // if (this->modelBin.find(".weights") != string::npos)
  // {
  //   this->net = readNetFromDarknet(this->modelTxt, this->modelBin);
  // }
  // Init classes
  string classesFile;
  if (!this->config->get("classesFile", classesFile))
  {
    cerr << "MatcherDnn: Missing classesFile" << endl;
    return;
  }
  this->loadClasses(classesFile);

  if (this->net.empty())
  {
    cerr << "MatcherDnn: DNN is empty" << endl;
    return;
  }
}

vector<Match> MatcherDnn::match(Camera *cam)
{
  // namedWindow("dnn", WINDOW_NORMAL);
  cam->getFrame().copyTo(this->frame);
  Size size = this->frame.size();
  // First, squarize roi to improve detection
  vector<cv::Rect> rects = squarizeRois(size, cam->getRois(), this->inputBlobSize.width);
  // Merge nearest roi
  rects = mergeRois(rects, 10);
  rects = squarizeRois(size, rects);
  vector<Match> matches;
  for (int i = 0; i < rects.size(); i++)
  {
    rectangle(this->frame, rects.at(i), cv::Scalar(255, 255, 0), 2);
    this->detectMatchesOnRoi(cam, rects.at(i), &matches);
  }
  //imshow("dnn", this->frame);
  return matches;
}

string loadClasses_trimValue(const string str)
{
  int first = str.find_first_not_of(" \t");
  if (first != string::npos)
  {
    int last = str.find_last_not_of(" \t");
    return str.substr(first, last - first + 1);
  }
  return "";
}

bool MatcherDnn::loadClasses(const string filename)
{
  ifstream file(filename.c_str());

  if (!file.good())
  {
    cerr << "MatcherDnn: Cannot read classes file: " << filename << endl;
    return false;
  }
  this->classes.clear();
  while (file.good() && !file.eof())
  {
    string line;
    getline(file, line);
    // split line into key and value
    if (!line.empty())
    {
      int pos = line.find_first_of(" \t");

      if (pos != string::npos)
      {
        string key = line.substr(0, pos);
        string value = loadClasses_trimValue(line.substr(pos + 1));

        if (!key.empty() && !value.empty())
        {
          this->classes[key] = value;
        }
      }
    }
  }
  return true;
}
