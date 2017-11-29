#ifndef MATCHER_H
#define MATCHER_H

#include <vector>
#include <map>
#include <string>
#include <opencv2/opencv.hpp>

typedef struct Match Match;

#include "config.hpp"
#include "camera.hpp"

using namespace std;

struct Match
{
  string name;
  float confidence;
  cv::Rect roi;
};

// Forward declarations
class Camera;

class Matcher
{
public:
  void setConfig(Config *config);
  virtual void init() = 0;
  virtual vector<Match> match(Camera *cam) = 0;

protected:
  Config *config;
};

#endif //MATCHER_H
