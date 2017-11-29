#include "camera-binding.hpp"

#include <string>
#include <vector>
#include <cmath>
#include <opencv2/opencv.hpp>
#include "camera.hpp"
#include "matcher.hpp"
#include "matcher-dnn.hpp"

struct CameraHandler
{
  Camera *cam;
};

struct CameraMatch
{
  Mat frame;
  Camera *cam;
  CameraMatchList *matches = NULL;
};

map<string, Matcher *> matchers;

using namespace cv;

EXTERNC CameraHandler *camera_create(const char *endpoint)
{
  CameraHandler *ch = new CameraHandler();
  ch->cam = new Camera(endpoint);
  return ch;
}

EXTERNC void camera_setPersistence(CameraHandler *ch, float value)
{
  ch->cam->setPersistence(value);
}

EXTERNC void camera_setMatcher(CameraHandler *ch, const char *name)
{
  ch->cam->setMatcher(matchers[string(name)]);
}

EXTERNC void camera_setMask(CameraHandler *ch, ZoneList *zones)
{
  vector<Rect> maskRects;
  for (ZoneList *node = zones; node != NULL; node = node->next)
  {
    ZoneItem *item = node->item;
    maskRects.push_back(Rect(item->left, item->top, item->width, item->height));
  }
  ch->cam->setMask(maskRects);
}

EXTERNC void camera_start(CameraHandler *ch)
{
  ch->cam->start();
}

EXTERNC void camera_stop(CameraHandler *ch)
{
  ch->cam->stop();
}

EXTERNC void camera_destroy(CameraHandler *ch)
{
  delete ch->cam;
  free(ch);
}

EXTERNC CameraMatch *cameraMatch_init(CameraHandler *ch)
{
  CameraMatch *cm = new CameraMatch();
  cm->cam = ch->cam;
  return cm;
}

EXTERNC bool cameraMatch_check(CameraMatch *cm)
{
  if (cm->cam->check())
  {
    cm->cam->getFrame().copyTo(cm->frame);
    vector<Match> matches = cm->cam->getMatches();

    CameraMatchList *node = NULL;
    cm->matches = NULL;
    // Get match as linked list
    for (int i = 0; i < matches.size(); i++)
    {
      Match m = matches.at(i);
      CameraMatchItem *item = new CameraMatchItem();
      item->label = new char[m.name.length() + 1];
      strcpy(item->label, m.name.c_str());
      item->confidence = m.confidence;
      item->left = m.roi.x;
      item->top = m.roi.y;
      item->width = m.roi.width;
      item->height = m.roi.height;
      node = cameraMatchList_append(node, item);
      if (cm->matches == NULL)
      {
        cm->matches = node;
      }
    }
    return true;
  }
  return false;
}

EXTERNC void matcher_init(const char *name, const char *matcherType, const char *configStr)
{
  if (matchers[name] == NULL)
  {
    Config *config = new Config();
    config->loadString(configStr);
    Matcher *matcher;
    if (strcmp(matcherType, "dnn") == 0)
    {
      matcher = new MatcherDnn();
    }
    else
    {
      return;
    }

    matcher->setConfig(config);
    matcher->init();
    matchers[name] = matcher;
  }
}

EXTERNC CameraMatchList *cameraMatch_getMatches(CameraMatch *cm)
{
  return cm->matches;
}

Scalar stringToScalar(char *hexValue)
{
  int r, g, b;
  sscanf(hexValue, "%02x%02x%02x", &b, &g, &r);
  return Scalar(r, g, b);
}

EXTERNC void cameraMatch_capture(CameraMatch *cm, const char *filename, ZoneList *zones)
{
  Mat frame(cm->frame);
  if (zones != NULL)
  {
    for (ZoneList *node = zones; node != NULL; node = node->next)
    {
      ZoneItem *item = node->item;
      Rect roi(item->left, item->top, item->width, item->height);
      Scalar color = stringToScalar(item->color);
      if (item->fillOpacity > 0)
      {
        Mat overlay;
        double alpha = (double)item->fillOpacity;
        frame.copyTo(overlay);
        rectangle(overlay, roi, color, -1);
        addWeighted(overlay, alpha, frame, 1 - alpha, 0, frame);
      }
      if (item->borderSize > 0)
      {
        rectangle(frame, roi, color, item->borderSize);
      }
      if (item->label != NULL)
      {
        String label = String(item->label);
        int baseLine = 0;
        Size labelSize = getTextSize(label, FONT_HERSHEY_SIMPLEX, 0.5, 1, &baseLine);
        rectangle(frame, Rect(Point(roi.x, roi.y - labelSize.height), Size(labelSize.width, labelSize.height + baseLine)),
                  color, CV_FILLED);
        putText(frame, label, Point(roi.x, roi.y),
                FONT_HERSHEY_SIMPLEX, 0.5, Scalar(0, 0, 0));
      }
    }
  }
  imwrite(filename, frame);
}

EXTERNC ImageSize cameraMatch_getImageSize(CameraMatch *cm)
{
  Size size = cm->frame.size();
  return ImageSize{size.width, size.height};
}

EXTERNC void cameraMatch_destroy(CameraMatch *cm)
{
  cm->cam = NULL;
  cameraMatchList_destroy(cm->matches, true);
  free(cm);
}

EXTERNC CameraMatchList *cameraMatchList_append(CameraMatchList *head, CameraMatchItem *cmi)
{
  CameraMatchList *node = new CameraMatchList();
  node->item = cmi;
  node->next = NULL;
  if (head == NULL)
  {
    return node;
  }
  if (head->next)
  {
    // Insert
    node->next = head->next;
  }
  head->next = node;
  return node;
}
EXTERNC void cameraMatchList_destroy(CameraMatchList *list, bool freeItems = false)
{
  CameraMatchList *node;
  while (list != NULL)
  {
    node = list;
    list = node->next;
    if (freeItems)
    {
      delete[] node->item->label;
      free(node->item);
    }
    free(node);
  }
}

EXTERNC ZoneItem *zoneItem_new()
{
  return new ZoneItem();
}

EXTERNC ZoneList *zoneList_append(ZoneList *head, ZoneItem *zi)
{
  ZoneList *node = new ZoneList();
  node->item = zi;
  node->next = NULL;
  if (head == NULL)
  {
    return node;
  }
  if (head->next)
  {
    // Insert
    node->next = head->next;
  }
  head->next = node;
  return node;
}
EXTERNC void zoneList_destroy(ZoneList *list, bool freeItems)
{
  ZoneList *node;
  while (list != NULL)
  {
    node = list;
    list = node->next;
    if (freeItems && node->item != NULL)
    {
      if (node->item->label != NULL)
        delete[] node->item->label;
      if (node->item->color != NULL)
        delete[] node->item->color;
      free(node->item);
    }
    free(node);
  }
}