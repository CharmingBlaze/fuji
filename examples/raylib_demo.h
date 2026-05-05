/*******************************************************************************************
*
*   raylib - A simple and easy-to-use library to enjoy videogames programming
*
*   Copyright (c) 2013-2023 Ramon Santamaria (@raysan5)
*
********************************************************************************************/

#ifndef RAYLIB_H
#define RAYLIB_H

#include <stdbool.h>
#include <stdint.h>

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
typedef bool bool;

typedef struct Vector2 {
    float x;
    float y;
} Vector2;

typedef struct Vector3 {
    float x;
    float y;
    float z;
} Vector3;

typedef struct Vector4 {
    float x;
    float y;
    float z;
    float w;
} Vector4;

typedef struct Color {
    unsigned char r;
    unsigned char g;
    unsigned char b;
    unsigned char a;
} Color;

typedef struct Rectangle {
    float x;
    float y;
    float width;
    float height;
} Rectangle;

typedef struct Image {
    void *data;
    int width;
    int height;
    int mipmaps;
    int format;
} Image;

typedef struct Texture2D {
    unsigned int id;
    int width;
    int height;
    int mipmaps;
    int format;
} Texture2D;

typedef struct Camera3D {
    Vector3 position;
    Vector3 target;
    Vector3 up;
    float fovy;
    int projection;
} Camera3D;

//----------------------------------------------------------------------------------
// Enumerations Definition
//----------------------------------------------------------------------------------
typedef enum {
    FLAG_VSYNC_HINT = 0x00000040,
    FLAG_FULLSCREEN_MODE = 0x00000002,
    FLAG_WINDOW_RESIZABLE = 0x00000004,
    FLAG_WINDOW_UNDECORATED = 0x00000008,
    FLAG_WINDOW_TRANSPARENT = 0x00000010,
    FLAG_WINDOW_HIDDEN = 0x00000080,
    FLAG_MSAA_4X_HINT = 0x00000020,
    FLAG_INTERLACED_HINT = 0x00010000
} ConfigFlags;

typedef enum {
    LOG_ALL = 0,
    LOG_TRACE,
    LOG_DEBUG,
    LOG_INFO,
    LOG_WARNING,
    LOG_ERROR,
    LOG_FATAL,
    LOG_NONE
} TraceLogLevel;

//----------------------------------------------------------------------------------
// Global Functions Declaration
//----------------------------------------------------------------------------------

// Window-related functions
void InitWindow(int width, int height, const char *title);
bool WindowShouldClose(void);
void CloseWindow(void);
bool IsWindowReady(void);
bool IsWindowFullscreen(void);
void ToggleFullscreen(void);
void SetWindowIcon(Image image);
void SetWindowTitle(const char *title);
void SetWindowPosition(int x, int y);
void SetWindowMonitor(int monitor);
void SetWindowMinSize(int width, int height);
void SetWindowSize(int width, int height);
void SetWindowOpacity(float opacity);
int GetScreenWidth(void);
int GetScreenHeight(void);
int GetMonitorCount(void);
int GetCurrentMonitor(void);
int GetMonitorWidth(int monitor);
int GetMonitorHeight(int monitor);
int GetMonitorRefreshRate(int monitor);
Vector2 GetWindowPosition(void);
Vector2 GetWindowScaleDPI(void);
const char *GetMonitorName(int monitor);

// Cursor-related functions
void ShowCursor(void);
void HideCursor(void);
bool IsCursorHidden(void);
void EnableCursor(void);
void DisableCursor(void);

// Drawing-related functions
void BeginDrawing(void);
void EndDrawing(void);
void BeginMode2D(Camera2D camera);
void EndMode2D(void);
void BeginMode3D(Camera3D camera);
void EndMode3D(void);
void BeginTextureMode(RenderTexture2D target);
void EndTextureMode(void);

// Basic shapes drawing functions
void DrawPixel(int posX, int posY, Color color);
void DrawPixelV(Vector2 position, Color color);
void DrawLine(int startPosX, int startPosY, int endPosX, int endPosY, Color color);
void DrawLineV(Vector2 startPos, Vector2 endPos, Color color);
void DrawLineEx(Vector2 startPos, Vector2 endPos, float thick, Color color);
void DrawLineStrip(Vector2 *points, int pointCount, Color color);
void DrawCircle(int centerX, int centerY, float radius, Color color);
void DrawCircleSector(Vector2 center, float radius, float startAngle, float endAngle, int segments, Color color);
void DrawCircleSectorLines(Vector2 center, float radius, float startAngle, float endAngle, int segments, Color color);
void DrawCircleGradient(int centerX, int centerY, float radius, Color color1, Color color2);
void DrawCircleV(Vector2 center, float radius, Color color);
void DrawCircleLines(int centerX, int centerY, float radius, Color color);
void DrawEllipse(int centerX, int centerY, float radiusH, float radiusV, Color color);
void DrawEllipseLines(int centerX, int centerY, float radiusH, float radiusV, Color color);
void DrawRing(Vector2 center, float innerRadius, float outerRadius, float startAngle, float endAngle, int segments, Color color);
void DrawRingLines(Vector2 center, float innerRadius, float outerRadius, float startAngle, float endAngle, int segments, Color color);
void DrawRectangle(int posX, int posY, int width, int height, Color color);
void DrawRectangleV(Vector2 position, Vector2 size, Color color);
void DrawRectangleRec(Rectangle rec, Color color);
void DrawRectanglePro(Rectangle rec, Vector2 origin, float rotation, Color color);
void DrawRectangleGradientV(int posX, int posY, int width, int height, Color color1, Color color2);
void DrawRectangleGradientH(int posX, int posY, int width, int height, Color color1, Color color2);
void DrawRectangleGradientEx(Rectangle rec, Color color1, Color color2, Color color3, Color color4);
void DrawRectangleLines(int posX, int posY, int width, int height, Color color);
void DrawRectangleLinesEx(Rectangle rec, float thick, Color color);
void DrawRectangleRounded(Rectangle rec, float roundness, int segments, Color color);
void DrawRectangleRoundedLines(Rectangle rec, float roundness, int segments, float thick, Color color);
void DrawTriangle(Vector2 v1, Vector2 v2, Vector2 v3, Color color);
void DrawTriangleLines(Vector2 v1, Vector2 v2, Vector2 v3, Color color);
void DrawTriangleFan(Vector2 *points, int pointCount, Color color);
void DrawPoly(Vector2 center, int sides, float radius, float rotation, Color color);
void DrawPolyLines(Vector2 center, int sides, float radius, float rotation, Color color);
void DrawPolyLinesEx(Vector2 center, int sides, float radius, float rotation, float lineThick, Color color);

// Basic shapes collision detection functions
bool CheckCollisionRecs(Rectangle rec1, Rectangle rec2);
bool CheckCollisionCircles(Vector2 center1, float radius1, Vector2 center2, float radius2);
bool CheckCollisionCircleRec(Vector2 center, float radius, Rectangle rec);
bool CheckCollisionPointRec(Vector2 point, Rectangle rec);
bool CheckCollisionPointCircle(Vector2 point, Vector2 center, float radius);
bool CheckCollisionPointTriangle(Vector2 point, Vector2 p1, Vector2 p2, Vector2 p3);
bool CheckCollisionLines(Vector2 startPos1, Vector2 endPos1, Vector2 startPos2, Vector2 endPos2, Vector2 *collisionPoint);
bool CheckCollisionPointPoly(Vector2 point, Vector2 *points, int pointCount);

// Texture loading/unloading functions
Texture2D LoadTexture(const char *fileName);
Image LoadImage(const char *fileName);
void UnloadTexture(Texture2D texture);
void UnloadImage(Image image);
Image LoadImageRaw(const char *fileName, int width, int height, int format, int headerSize);
Image LoadImageSVG(const char *fileNameOrString, int width, int height);
Texture2D LoadTextureFromImage(Image image);
Image LoadImageFromTexture(Texture2D texture);
bool ExportImage(Image image, const char *fileName);
bool ExportImageAsCode(Image image, const char *fileName);

// Image manipulation functions
Image ImageCopy(Image image);
Image ImageText(const char *text, int fontSize, Color color);
Image ImageTextEx(Font font, const char *text, float fontSize, float spacing, Color tint);
void ImageFormat(Image *image, int newFormat);
void ImageToPOT(Image *image, Color fill);
void ImageCrop(Image *image, Rectangle crop);
void ImageAlphaCrop(Image *image, float threshold);
void ImageAlphaClear(Image *image, Color color, float threshold);
void ImageAlphaMask(Image *image, Image alphaMask);
void ImageBlurGaussian(Image *image, int blurSize);
void ImageResize(Image *image, int newWidth, int newHeight);
void ImageResizeNN(Image *image, int newWidth, int newHeight);
void ImageResizeCanvas(Image *image, int newWidth, int newHeight, int offsetX, int offsetY, Color fill);
void ImageMipmaps(Image *image);
void ImageDither(Image *image, int rBpp, int gBpp, int bBpp, int aBpp);
void ImageFlipVertical(Image *image);
void ImageFlipHorizontal(Image *image);
void ImageRotateCW(Image *image);
void ImageRotateCCW(Image *image);
void ImageColorTint(Image *image, Color color);
void ImageColorInvert(Image *image);
void ImageColorGrayscale(Image *image);
void ImageColorContrast(Image *image, float contrast);
void ImageColorBrightness(Image *image, int brightness);
void ImageColorReplace(Image *image, Color color, Color replace);

// Color/pixel formats helper functions
int GetPixelDataSize(int width, int height, int format);
Color GetColor(unsigned int hexValue);
Color GetPixelColor(void *srcPtr, int format);
void SetPixelColor(void *dstPtr, Color color, int format);

// Font loading/unloading functions
Font GetFontDefault(void);
Font LoadFont(const char *fileName);
Font LoadFontEx(const char *fileName, int fontSize, int *fontChars, int glyphCount);
Font LoadFontFromImage(Image image, Color key, int firstChar);
GlyphInfo *LoadFontData(const char *fileName, int fontSize, int *fontChars, int glyphCount, int type);
Image GenImageFontAtlas(const GlyphInfo *glyphs, Rectangle **glyphRecs, int glyphCount, int fontSize, int padding, int packMethod);
void UnloadFont(Font font);

// Text drawing functions
void DrawFPS(int posX, int posY);
void DrawText(const char *text, int posX, int posY, int fontSize, Color color);
void DrawTextEx(Font font, const char *text, Vector2 position, float fontSize, float spacing, Color tint);
void DrawTextPro(Font font, const char *text, Vector2 position, Vector2 origin, float rotation, float fontSize, float spacing, Color tint);
void DrawTextCodepoint(Font font, int codepoint, Vector2 position, float fontSize, Color tint);
void DrawTextCodepoints(Font font, const int *codepoints, int codepointCount, Vector2 position, float fontSize, float spacing, Color tint);

// Text misc. functions
int MeasureText(const char *text, int fontSize);
Vector2 MeasureTextEx(Font font, const char *text, float fontSize, float spacing);
int GetGlyphIndex(Font font, int codepoint);
GlyphInfo GetGlyphInfo(Font font, int codepoint);
Rectangle GetGlyphAtlasRec(Font font, int codepoint);

// Camera system functions
Camera2D GetCamera2D(void);
Camera3D GetCamera3D(void);
void SetCamera2D(Camera2D camera);
void SetCamera3D(Camera3D camera);
void SetCameraPanControl(int panKey);
void SetCameraAltControl(int altKey);
void SetCameraSmoothZoomControl(int szKey);
void SetCameraMoveControls(int frontKey, int backKey, int rightKey, int leftKey, int upKey, int downKey);

// Misc. functions
void SetTraceLogLevel(int logType);
void SetTraceLogExit(int logType);
void TraceLog(int logType, const char *text, ...);
void TakeScreenshot(const char *fileName);
int GetRandomValue(int min, int max);

// Files management functions
bool FileExists(const char *fileName);
bool DirectoryExists(const char *dirPath);
bool IsFileExtension(const char *fileName, const char *extension);
bool IsPathFile(const char *path);
const char *GetFileExtension(const char *fileName);
const char *GetFileName(const char *filePath);
const char *GetFileNameWithoutExt(const char *filePath);
const char *GetDirectoryPath(const char *filePath);
const char *GetPrevDirectoryPath(const char *dirPath);
const char *GetWorkingDirectory(void);
const char *GetApplicationDirectory(void);
bool ChangeDirectory(const char *dir);
bool IsFileDropped(void);
char **GetDroppedFiles(int *count);
void ClearDroppedFiles(void);

// Persistent storage management
void SaveStorageValue(unsigned int position, int value);
int LoadStorageValue(unsigned int position);

// Input-related functions: keyboard
bool IsKeyPressed(int key);
bool IsKeyDown(int key);
bool IsKeyReleased(int key);
bool IsKeyUp(int key);
int GetKeyPressed(void);
void SetExitKey(int key);

// Input-related functions: gamepads
bool IsGamepadAvailable(int gamepad);
bool IsGamepadName(int gamepad, const char *name);
const char *GetGamepadName(int gamepad);
bool IsGamepadButtonPressed(int gamepad, int button);
bool IsGamepadButtonDown(int gamepad, int button);
bool IsGamepadButtonReleased(int gamepad, int button);
bool IsGamepadButtonUp(int gamepad, int button);
int GetGamepadButtonPressed(void);
int GetGamepadAxisCount(int gamepad);
float GetGamepadAxisMovement(int gamepad, int axis);

// Input-related functions: mouse
bool IsMouseButtonPressed(int button);
bool IsMouseButtonDown(int button);
bool IsMouseButtonReleased(int button);
bool IsMouseButtonUp(int button);
int GetMouseX(void);
int GetMouseY(void);
Vector2 GetMousePosition(void);
Vector2 GetMouseDelta(void);
void SetMousePosition(int x, int y);
void SetMouseOffset(int offsetX, int offsetY);
void SetMouseScale(float scaleX, float scaleY);
float GetMouseWheelMove(void);
int GetMouseCursor(void);
void SetMouseCursor(int cursor);

// Touch-related functions
int GetTouchX(void);
int GetTouchY(void);
Vector2 GetTouchPosition(int index);
int GetTouchPointCount(void);
int GetTouchPointId(int index);

// Gestures-related functions
void SetGesturesEnabled(unsigned int flags);
bool IsGestureDetected(unsigned int gesture);
int GetGestureDetected(void);
void GetGestureDragVector(int *vectorX, int *vectorY);
void GetGestureDragAngle(float *angle);
float GetGesturePinchVector(void);
float GetGesturePinchAngle(void);

// Camera3D type (depends on Vector3 and Camera3D)
typedef struct Camera2D {
    Vector2 offset;
    Vector2 target;
    float rotation;
    float zoom;
} Camera2D;

typedef struct RenderTexture2D {
    unsigned int id;
    Texture2D texture;
    Texture2D depth;
    bool depthTexture;
} RenderTexture2D;

typedef struct Font {
    int baseSize;
    int glyphCount;
    int glyphPadding;
    Texture2D texture;
    Rectangle *recs;
    GlyphInfo *glyphs;
} Font;

typedef struct GlyphInfo {
    int value;
    int offsetX;
    int offsetY;
    int advanceX;
    Image image;
} GlyphInfo;

#endif // RAYLIB_H
