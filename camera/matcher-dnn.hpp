#ifndef MOTION_MACTHER_DNN_H
#define MOTION_MACTHER_DNN_H

#include <vector>
#include <map>
#include <string>
#include <opencv2/dnn.hpp>

#include "matcher.hpp"
#include "camera.hpp"

using namespace cv;
using namespace cv::dnn;

class MatcherDnn : public Matcher
{
public:
  MatcherDnn() {
    this->roiPadding = 10;
    this->minConfidence = 0.1f;
    this->modelTxt = "";
    this->modelBin = "";
    this->netInput = "data";
    this->netOutput = "detection_out";
    this->inScaleFactor = 1.0f;
    this->meanVal = 127.5;
    this->inputBlobSize = Size(300, 300);
  };
  void init();
  vector<Match> match(Camera *cam);

private:
  void detectMatchesOnRoi(Camera *cam, Rect roi, vector<Match> *matches);
  bool loadClasses(const string filename);

  map<string, string> classes;

  Net net;
  Mat frame;
  // Settings
  int roiPadding;
  float minConfidence;
  string modelTxt;
  string modelBin;
  string netInput;
  string netOutput;
  float inScaleFactor;
  float meanVal;
  Size inputBlobSize;
};

#endif //MOTION_MACTHER_DNN_H
