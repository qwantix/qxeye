
#include <iostream>
#include <unistd.h>
#include <chrono>
#include <opencv2/opencv.hpp>
#include "opencv2/videoio.hpp"
#include "opencv2/imgproc.hpp"
// #include "opencv2/objdetect.hpp"

#include "camera.hpp"
#include "config.hpp"
#include "matcher-dnn.hpp"
#include "util.hpp"

using namespace std;
using namespace cv;

Camera::Camera(const char *endpoint)
{
  cout << "camera: Open endpoint " << endpoint << endl;
  // Init video capture
  this->cap = new VideoCapture(endpoint);
  this->cap->set(CV_CAP_PROP_BUFFERSIZE, 1);
  if (!this->cap->isOpened())
  {
    cerr << "camera: Cannot open the video file" << endl;
  }
}

Camera::~Camera()
{
  cout << "camera: Destroy camera" << endl;
  this->reading = false;
  if (!this->cap->isOpened())
  {
    this->cap->release();
  }
  delete this->cap;
}

void Camera::setPersistence(float value)
{
  this->persistence = value;
}

void Camera::setMatcher(Matcher *matcher)
{
  this->matcher = matcher;
}

void Camera::setMask(vector<Rect> mask)
{
  this->mask = mask;
}

Mat Camera::getFrame()
{
  return this->frame;
}

Mat Camera::getFrame(Rect r)
{
  return this->frame(r);
}

vector<Rect> Camera::getRois()
{
  return this->rois;
}

vector<Match> Camera::getMatches()
{
  return this->matches;
}

void Camera::start()
{
  this->reading = true;
  while (this->reading)
  {
    if (this->cap->grab())
    {
      this->frameInvalid = true;
    }
    usleep(10000);
  }
}

void Camera::stop()
{
  this->reading = false;
}

bool Camera::check()
{
  if (this->matcher == NULL)
  {
    return false;
  }
  if (!this->fetchFrame())
  {
    return false;
  }
  if (!this->fetchMotionRois())
  {
    return false;
  }
  if (this->rois.size() == 0)
  {
    return false;
  }
  this->matches = this->matcher->match(this);
  return this->matches.size() > 0;
}

bool Camera::fetchFrame()
{
  if (!this->frameInvalid)
  {
    return false;
  }
  if (!this->cap->isOpened())
  {
    cerr << "camera: Camera not open" << endl;
    return false;
  }

  bool b = this->cap->retrieve(this->frame);
  this->frameInvalid = false;
  return b;
}

bool Camera::fetchMotionRois()
{
  int detectW = 150;
  Mat frame;
  Size size = this->frame.size();
  Size resizedSize;
  float ratio = (float)size.width / (float)detectW;
  resizedSize.width = detectW;
  resizedSize.height = ((float)size.height) / ratio;

  if (this->maskMat.empty())
  {
    this->initMaskMat(resizedSize);
  }

  resize(this->frame, frame, resizedSize);
  cvtColor(frame, frame, cv::COLOR_BGR2GRAY);
  equalizeHist(frame, frame);
  bitwise_and(this->maskMat, frame, frame);

  // Push frame to persitence
  float alpha = this->persistence / 1000;
  float beta = 1 - alpha;
  if (this->persistentFrame.rows == 0)
  {
    frame.copyTo(this->persistentFrame);
  }
  else
  {
    addWeighted(frame, alpha, this->persistentFrame, beta, 0.0, this->persistentFrame);
  }
  //imshow( "Persistence", this->persistentFrame );
  absdiff(frame, this->persistentFrame, frame);
  //imshow("Diff", frame); // Disabled
  threshold(frame, frame, 100, 255, THRESH_BINARY);
  //adaptiveThreshold(frame, thres,255,ADAPTIVE_THRESH_GAUSSIAN_C, CV_THRESH_BINARY,15,5);
  //imshow("S1", frame);
  blur(frame, frame, Size(3, 3));
  //imshow("Blur", frame);
  threshold(frame, frame, 100, 255, THRESH_BINARY);
  int morphElem = 0; // 0: Rect - 1: Cross - 2: Ellipse
  int morphSize = 5; // Kernel size:\n 2n +1
  Mat kernel = getStructuringElement(morphElem, Size(2 * morphSize + 1, 2 * morphSize + 1), Point(morphSize, morphSize));
  morphologyEx(frame, frame, MORPH_CLOSE, kernel);
  //imshow( "S2", frame ); // Disabled

  // adaptative pause
  int pauseMs = this->rois.size() > 0 ? 100 : 300;
  usleep(pauseMs * 1000);
  //waitKey(pauseMs);

  std::vector<std::vector<Point>>
      contours;
  std::vector<Vec4i> hierarchy;

  findContours(frame, contours, hierarchy, RETR_EXTERNAL, CHAIN_APPROX_SIMPLE); // retrieves external contours

  if (contours.size() == 0)
  {
    return false;
  }

  cvtColor(frame, frame, CV_GRAY2RGB);

  std::vector<Rect> rects;
  Rect r(0, 0, 0, 0);
  bool first = true;

  this->rois.clear();
  for (int i = 0; i < contours.size(); i++)
  {
    r = boundingRect(contours.at(i));
    if (r.width > 3 && r.height > 3)
    {
      rectangle(frame, r, Scalar(0, 0, 255), 1);
      rects.push_back(r);
      Rect roiRescaled(r);
      roiRescaled.x = max((int)(roiRescaled.x * ratio), 0);
      roiRescaled.y = max((int)(roiRescaled.y * ratio), 0);
      roiRescaled.width = min((int)(roiRescaled.width * ratio), (int)(size.width - roiRescaled.x));
      roiRescaled.height = min((int)(roiRescaled.height * ratio), (int)(size.height - roiRescaled.y));

      this->rois.push_back(roiRescaled);
      first = false;
    }
  }
  //imshow("Out", frame);
  return true;
}

void Camera::initMaskMat(Size size)
{
  Mat m(size, CV_8UC1, Scalar(255, 255, 255));
  for (int i = 0; i < this->mask.size(); i++)
  {
    Rect r(this->mask.at(i));
    r.x = size.width * r.x / 100;
    r.y = size.height * r.y / 100;
    r.width = size.width * r.width / 100;
    r.height = size.height * r.height / 100;
    rectangle(m, r, Scalar(0, 0, 0), CV_FILLED);
  }
  this->maskMat = m;
}