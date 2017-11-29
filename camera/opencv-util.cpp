
#include <iostream>
#include <opencv2/opencv.hpp>
#include "opencv-util.hpp"
#include "util.hpp"

using namespace std;
using namespace cv;

Rect adjustRectToSize(Rect rect, Size size)
{
  Rect out(rect);
  if (size.width > rect.width)
  {
    out.x -= (int)((size.width - rect.width) / 2);
    out.width = size.width;
  }
  if (size.height > rect.height)
  {
    out.y -= (int)((size.height - rect.height) / 2);
    out.height = size.height;
  }
  return out;
}

std::vector<cv::Rect> squarizeRois(Size imgSize, std::vector<cv::Rect> rois, int minSize, int padding)
{
  int roisLen = rois.size();
  std::vector<cv::Rect> out;
  for (int i = 0; i < roisLen; i++)
  {
    Rect originRoi = rois.at(i);
    Rect roi(originRoi);

    // Convert roi as square
    // get max size
    int maxSize = max(roi.width, roi.height);
    maxSize = max(maxSize - padding, minSize - padding);
    // Add margin
    maxSize += padding;
    roi = adjustRectToSize(roi, Size(maxSize, maxSize));

    // OK, now adjust square into imgSize bounds
    if (roi.height > imgSize.height)
    {
      roi.height = imgSize.height;
      roi.width = imgSize.width;
    }
    else if (roi.width > imgSize.width)
    {
      roi.height = imgSize.height;
      roi.width = imgSize.width;
    }

    roi.x = max(0, roi.x);
    roi.y = max(0, roi.y);
    if ((roi.x + roi.width) > imgSize.width)
    {
      roi.width = imgSize.width - roi.x;
    }
    if ((roi.y + roi.height) > imgSize.height)
    {
      roi.height = imgSize.height - roi.y;
    }

    if (false && roi.height != roi.width)
    {
      Rect roi2;
      // create sliding window when rect require 2 squares
      if (roi.height > roi.width)
      { // portrait
        roi2.height = roi2.width = roi.width;
        roi2.x = roi.x;
        roi2.y = roi.height - roi.width;
        roi.height = roi.width;
      }
      else
      { // landscape
        roi2.width = roi2.height = roi.height;
        roi2.y = roi.y;
        roi2.x = roi.width - roi.height;
        roi.width = roi.height;
      }
      out.push_back(roi2);
    }

    out.push_back(roi);
  }
  return out;
}

bool isNear(Rect a, Rect b, int distance)
{
  return b.x > a.x ? b.x - a.x <= a.width + distance : a.x - b.x <= b.width + distance && b.y > a.y ? b.y - a.y <= a.height + distance : a.y - b.y <= b.height + distance;
}

std::vector<cv::Rect> mergeRois(std::vector<cv::Rect> rois, int distance)
{
  std::vector<cv::Rect> out = rois;
  std::vector<cv::Rect> subset;
  bool merged = false;
  do
  {
    Rect a, b;
    subset = out;
    out.clear();
    int len = subset.size();
    merged = false;
    for (int i = 0; i < len && !merged; i++)
    {
      a = Rect(subset.at(i));
      for (int j = i + 1; j < len; j++)
      {
        b = subset.at(j);
        if (isNear(a, b, distance))
        {
          a.width = max(a.x + a.width, b.x + b.width) - min(a.x, b.x);
          a.height = max(a.y + a.height, b.y + b.height) - min(a.y, b.y);
          a.x = min(a.x, b.x);
          a.y = min(a.y, b.y);
          merged = true;
        }
      }
      out.push_back(a);
    }
  } while (merged);
  return out;
}
