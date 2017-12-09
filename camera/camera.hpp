#ifndef CAMERA_H
#define CAMERA_H

#include <string>
#include <vector>
#include <opencv2/opencv.hpp>

#include "matcher.hpp"

using namespace std;
using namespace cv;

// Forward declarations
class Matcher;

class Camera
{
public:
  Camera(const char *endpoint);
  virtual ~Camera();

  void setPersistence(float value);
  void setMatcher(Matcher *matcher);
  void setMask(vector<Rect> mask);
  Mat getFrame();
  Mat getFrame(Rect r);
  vector<Rect> getRois();
  vector<Match> getMatches();

  void start();
  void stop();
  bool check();

private:
  bool fetchFrame();
  bool fetchMotionRois();
  void initMaskMat(Size size);

  VideoCapture *cap;
  vector<Rect> mask;
  Mat maskMat;
  Matcher *matcher;
  Mat frame;
  vector<Rect> rois;
  vector<Match> matches;
  // Flags
  bool frameInvalid;
  bool reading;
  Mat persistentFrame;
  // Settings
  float persistence;
};

#endif //CAMERA_H
