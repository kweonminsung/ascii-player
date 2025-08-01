#ifndef _OPENCV3_HIGHGUI_H_
#define _OPENCV3_HIGHGUI_H_

#ifdef __cplusplus
#include <opencv2/opencv.hpp>
extern "C" {
#endif

#include "core.h"

typedef void(*mouse_callback) (int event, int x, int y, int flags, void *userdata);

/* typedef struct mouse_callback_userdata {
    char* winname;
    void* userdata;
} mouse_callback_userdata; */

void Window_SetMouseCallback(char* winname, mouse_callback on_mouse);

// Window
void Window_New(const char* winname, int flags);
void Window_Close(const char* winname);
OpenCVResult Window_IMShow(const char* winname, Mat mat);
double Window_GetProperty(const char* winname, int flag);
OpenCVResult Window_SetProperty(const char* winname, int flag, double value);
OpenCVResult Window_SetTitle(const char* winname, const char* title);
int Window_WaitKey(int);
int Window_WaitKeyEx(int);
int Window_PollKey(void);
OpenCVResult Window_Move(const char* winname, int x, int y);
OpenCVResult Window_Resize(const char* winname, int width, int height);
struct Rect Window_SelectROI(const char* winname, Mat img);
struct Rects Window_SelectROIs(const char* winname, Mat img);

// Trackbar
void Trackbar_Create(const char* winname, const char* trackname, int max);
void Trackbar_CreateWithValue(const char* winname, const char* trackname, int* value, int max);
int Trackbar_GetPos(const char* winname, const char* trackname);
void Trackbar_SetPos(const char* winname, const char* trackname, int pos);
void Trackbar_SetMin(const char* winname, const char* trackname, int pos);
void Trackbar_SetMax(const char* winname, const char* trackname, int pos);

#ifdef __cplusplus
}
#endif

#endif //_OPENCV3_HIGHGUI_H_
