#ifndef OPENCV_UTIL_H
#define OPENCV_UTIL_H

#include <vector>
#include <opencv2/opencv.hpp>

using namespace cv;
using namespace std;

Rect adjustRectToSize(Rect rect, Size size);

vector<Rect> squarizeRois(Size imgSize, vector<Rect> rois, int minSize = 0, int padding = 10);

bool isNear(Rect a, Rect b, int distance);

vector<Rect> mergeRois(vector<Rect> rois, int distance);

#endif //OPENCV_UTIL_H
