#ifndef CAMERA_C_H
#define CAMERA_C_H

#ifdef __cplusplus
#define EXTERNC extern "C"
#else
#define EXTERNC
#endif

typedef struct CameraHandler CameraHandler;
typedef struct CameraMatch CameraMatch;
typedef struct CameraMatchList CameraMatchList;
typedef struct CameraMatchItem CameraMatchItem;

typedef struct ZoneList ZoneList;
typedef struct ZoneItem ZoneItem;

typedef struct ImageSize ImageSize;

struct CameraMatchItem
{
  char *label; // Here const char* cause pointer bug with go binding, so we use char*
  float confidence;
  int top, left;
  int width, height;
};

struct CameraMatchList
{
  CameraMatchItem *item;
  CameraMatchList *next;
};

struct ZoneItem
{
  char *label; // Here const char* cause pointer bug with go binding, so we use char*
  int top, left;
  int width, height;
  char *color;
  int borderSize;
  float fillOpacity;
};

struct ZoneList
{
  ZoneItem *item;
  ZoneList *next;
};

struct ImageSize
{
  int width;
  int height;
};

EXTERNC CameraHandler *camera_create(const char *endpoint);
EXTERNC void camera_setPersistence(CameraHandler *ch, float value);
EXTERNC void camera_setMatcher(CameraHandler *ch, const char *filename);
EXTERNC void camera_setMask(CameraHandler *ch, ZoneList *zl);
EXTERNC void camera_start(CameraHandler *ch);
EXTERNC void camera_stop(CameraHandler *ch);
EXTERNC void camera_destroy(CameraHandler *ch);

EXTERNC void matcher_init(const char *name, const char *matcherType, const char *configStr);

EXTERNC CameraMatch *cameraMatch_init(CameraHandler *ch);
EXTERNC bool cameraMatch_check(CameraMatch *cm);
EXTERNC void cameraMatch_capture(CameraMatch *cm, const char *filename, ZoneList *zones);
EXTERNC CameraMatchList *cameraMatch_getMatches(CameraMatch *cm);
EXTERNC ImageSize cameraMatch_getImageSize(CameraMatch *cm);
EXTERNC void cameraMatch_destroy(CameraMatch *cm);

EXTERNC CameraMatchList *cameraMatchList_append(CameraMatchList *cl, CameraMatchItem *cm);
EXTERNC void cameraMatchList_destroy(CameraMatchList *cl, bool freeItems);

EXTERNC ZoneItem *zoneItem_new();
EXTERNC ZoneList *zoneList_append(ZoneList *zl, ZoneItem *item);
EXTERNC void zoneList_destroy(ZoneList *zl, bool freeItems);

#endif //CAMERA_C_H
