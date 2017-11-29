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
  void init();
  vector<Match> match(Camera *cam);

private:
  void detectMatchesOnRoi(Camera *cam, Rect roi, vector<Match> *matches);
  bool loadClasses(const string filename);

  map<string, string> classes;

  Net net;
  Mat frame;
  // Settings
  int roiPadding = 10;
  float minConfidence = 0.1f;
  string modelTxt = "";
  string modelBin = "";
  string netInput = "data";
  string netOutput = "detection_out";
  float inScaleFactor = 1.0f;
  float meanVal = 127.5;
  Size inputBlobSize = Size(300, 300);
};

#endif //MOTION_MACTHER_DNN_H
