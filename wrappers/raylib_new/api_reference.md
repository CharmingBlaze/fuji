# raylib - API Reference

## Contents

- [Functions](#functions)
- [Structs](#structs)
- [Enums](#enums)
- [Macros](#macros)

---

## Functions

### void

```c
typedef void()
```

**Fuji usage**

```fuji
let result = void();
```

---

### bool

```c
typedef bool()
```

**Fuji usage**

```fuji
let result = bool();
```

---

### bool

```c
typedef bool()
```

**Fuji usage**

```fuji
let result = bool();
```

---

### InitWindow

```c
RLAPI void InitWindow(int width, int height, const char *title)
```

**Fuji usage**

```fuji
let result = InitWindow(width, height, *title);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `*title` | `const char` |

---

### CloseWindow

```c
RLAPI void CloseWindow()
```

**Fuji usage**

```fuji
let result = CloseWindow();
```

---

### WindowShouldClose

```c
RLAPI bool WindowShouldClose()
```

**Fuji usage**

```fuji
let result = WindowShouldClose();
```

---

### IsWindowReady

```c
RLAPI bool IsWindowReady()
```

**Fuji usage**

```fuji
let result = IsWindowReady();
```

---

### IsWindowFullscreen

```c
RLAPI bool IsWindowFullscreen()
```

**Fuji usage**

```fuji
let result = IsWindowFullscreen();
```

---

### IsWindowHidden

```c
RLAPI bool IsWindowHidden()
```

**Fuji usage**

```fuji
let result = IsWindowHidden();
```

---

### IsWindowMinimized

```c
RLAPI bool IsWindowMinimized()
```

**Fuji usage**

```fuji
let result = IsWindowMinimized();
```

---

### IsWindowMaximized

```c
RLAPI bool IsWindowMaximized()
```

**Fuji usage**

```fuji
let result = IsWindowMaximized();
```

---

### IsWindowFocused

```c
RLAPI bool IsWindowFocused()
```

**Fuji usage**

```fuji
let result = IsWindowFocused();
```

---

### IsWindowResized

```c
RLAPI bool IsWindowResized()
```

**Fuji usage**

```fuji
let result = IsWindowResized();
```

---

### IsWindowState

```c
RLAPI bool IsWindowState(unsigned int flag)
```

**Fuji usage**

```fuji
let result = IsWindowState(flag);
```

| Parameter | Type |
|-----------|------|
| `flag` | `unsigned int` |

---

### SetWindowState

```c
RLAPI void SetWindowState(unsigned int flags)
```

**Fuji usage**

```fuji
let result = SetWindowState(flags);
```

| Parameter | Type |
|-----------|------|
| `flags` | `unsigned int` |

---

### ClearWindowState

```c
RLAPI void ClearWindowState(unsigned int flags)
```

**Fuji usage**

```fuji
let result = ClearWindowState(flags);
```

| Parameter | Type |
|-----------|------|
| `flags` | `unsigned int` |

---

### ToggleFullscreen

```c
RLAPI void ToggleFullscreen()
```

**Fuji usage**

```fuji
let result = ToggleFullscreen();
```

---

### ToggleBorderlessWindowed

```c
RLAPI void ToggleBorderlessWindowed()
```

**Fuji usage**

```fuji
let result = ToggleBorderlessWindowed();
```

---

### MaximizeWindow

```c
RLAPI void MaximizeWindow()
```

**Fuji usage**

```fuji
let result = MaximizeWindow();
```

---

### MinimizeWindow

```c
RLAPI void MinimizeWindow()
```

**Fuji usage**

```fuji
let result = MinimizeWindow();
```

---

### RestoreWindow

```c
RLAPI void RestoreWindow()
```

**Fuji usage**

```fuji
let result = RestoreWindow();
```

---

### SetWindowIcon

```c
RLAPI void SetWindowIcon(Image image)
```

**Fuji usage**

```fuji
let result = SetWindowIcon(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### SetWindowIcons

```c
RLAPI void SetWindowIcons(Image *images, int count)
```

**Fuji usage**

```fuji
let result = SetWindowIcons(*images, count);
```

| Parameter | Type |
|-----------|------|
| `*images` | `Image` |
| `count` | `int` |

---

### SetWindowTitle

```c
RLAPI void SetWindowTitle(const char *title)
```

**Fuji usage**

```fuji
let result = SetWindowTitle(*title);
```

| Parameter | Type |
|-----------|------|
| `*title` | `const char` |

---

### SetWindowPosition

```c
RLAPI void SetWindowPosition(int x, int y)
```

**Fuji usage**

```fuji
let result = SetWindowPosition(x, y);
```

| Parameter | Type |
|-----------|------|
| `x` | `int` |
| `y` | `int` |

---

### SetWindowMonitor

```c
RLAPI void SetWindowMonitor(int monitor)
```

**Fuji usage**

```fuji
let result = SetWindowMonitor(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### SetWindowMinSize

```c
RLAPI void SetWindowMinSize(int width, int height)
```

**Fuji usage**

```fuji
let result = SetWindowMinSize(width, height);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |

---

### SetWindowMaxSize

```c
RLAPI void SetWindowMaxSize(int width, int height)
```

**Fuji usage**

```fuji
let result = SetWindowMaxSize(width, height);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |

---

### SetWindowSize

```c
RLAPI void SetWindowSize(int width, int height)
```

**Fuji usage**

```fuji
let result = SetWindowSize(width, height);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |

---

### SetWindowOpacity

```c
RLAPI void SetWindowOpacity(float opacity)
```

**Fuji usage**

```fuji
let result = SetWindowOpacity(opacity);
```

| Parameter | Type |
|-----------|------|
| `opacity` | `float` |

---

### SetWindowFocused

```c
RLAPI void SetWindowFocused()
```

**Fuji usage**

```fuji
let result = SetWindowFocused();
```

---

### GetScreenWidth

```c
RLAPI int GetScreenWidth()
```

**Fuji usage**

```fuji
let result = GetScreenWidth();
```

---

### GetScreenHeight

```c
RLAPI int GetScreenHeight()
```

**Fuji usage**

```fuji
let result = GetScreenHeight();
```

---

### GetRenderWidth

```c
RLAPI int GetRenderWidth()
```

**Fuji usage**

```fuji
let result = GetRenderWidth();
```

---

### GetRenderHeight

```c
RLAPI int GetRenderHeight()
```

**Fuji usage**

```fuji
let result = GetRenderHeight();
```

---

### GetMonitorCount

```c
RLAPI int GetMonitorCount()
```

**Fuji usage**

```fuji
let result = GetMonitorCount();
```

---

### GetCurrentMonitor

```c
RLAPI int GetCurrentMonitor()
```

**Fuji usage**

```fuji
let result = GetCurrentMonitor();
```

---

### GetMonitorPosition

```c
RLAPI Vector2 GetMonitorPosition(int monitor)
```

**Fuji usage**

```fuji
let result = GetMonitorPosition(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorWidth

```c
RLAPI int GetMonitorWidth(int monitor)
```

**Fuji usage**

```fuji
let result = GetMonitorWidth(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorHeight

```c
RLAPI int GetMonitorHeight(int monitor)
```

**Fuji usage**

```fuji
let result = GetMonitorHeight(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorPhysicalWidth

```c
RLAPI int GetMonitorPhysicalWidth(int monitor)
```

**Fuji usage**

```fuji
let result = GetMonitorPhysicalWidth(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorPhysicalHeight

```c
RLAPI int GetMonitorPhysicalHeight(int monitor)
```

**Fuji usage**

```fuji
let result = GetMonitorPhysicalHeight(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetMonitorRefreshRate

```c
RLAPI int GetMonitorRefreshRate(int monitor)
```

**Fuji usage**

```fuji
let result = GetMonitorRefreshRate(monitor);
```

| Parameter | Type |
|-----------|------|
| `monitor` | `int` |

---

### GetWindowPosition

```c
RLAPI Vector2 GetWindowPosition()
```

**Fuji usage**

```fuji
let result = GetWindowPosition();
```

---

### GetWindowScaleDPI

```c
RLAPI Vector2 GetWindowScaleDPI()
```

**Fuji usage**

```fuji
let result = GetWindowScaleDPI();
```

---

### SetClipboardText

```c
RLAPI void SetClipboardText(const char *text)
```

**Fuji usage**

```fuji
let result = SetClipboardText(*text);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |

---

### GetClipboardImage

```c
RLAPI Image GetClipboardImage()
```

**Fuji usage**

```fuji
let result = GetClipboardImage();
```

---

### EnableEventWaiting

```c
RLAPI void EnableEventWaiting()
```

**Fuji usage**

```fuji
let result = EnableEventWaiting();
```

---

### DisableEventWaiting

```c
RLAPI void DisableEventWaiting()
```

**Fuji usage**

```fuji
let result = DisableEventWaiting();
```

---

### ShowCursor

```c
RLAPI void ShowCursor()
```

**Fuji usage**

```fuji
let result = ShowCursor();
```

---

### HideCursor

```c
RLAPI void HideCursor()
```

**Fuji usage**

```fuji
let result = HideCursor();
```

---

### IsCursorHidden

```c
RLAPI bool IsCursorHidden()
```

**Fuji usage**

```fuji
let result = IsCursorHidden();
```

---

### EnableCursor

```c
RLAPI void EnableCursor()
```

**Fuji usage**

```fuji
let result = EnableCursor();
```

---

### DisableCursor

```c
RLAPI void DisableCursor()
```

**Fuji usage**

```fuji
let result = DisableCursor();
```

---

### IsCursorOnScreen

```c
RLAPI bool IsCursorOnScreen()
```

**Fuji usage**

```fuji
let result = IsCursorOnScreen();
```

---

### ClearBackground

```c
RLAPI void ClearBackground(Color color)
```

**Fuji usage**

```fuji
let result = ClearBackground(color);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |

---

### BeginDrawing

```c
RLAPI void BeginDrawing()
```

**Fuji usage**

```fuji
let result = BeginDrawing();
```

---

### EndDrawing

```c
RLAPI void EndDrawing()
```

**Fuji usage**

```fuji
let result = EndDrawing();
```

---

### BeginMode2D

```c
RLAPI void BeginMode2D(Camera2D camera)
```

**Fuji usage**

```fuji
let result = BeginMode2D(camera);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera2D` |

---

### EndMode2D

```c
RLAPI void EndMode2D()
```

**Fuji usage**

```fuji
let result = EndMode2D();
```

---

### BeginMode3D

```c
RLAPI void BeginMode3D(Camera3D camera)
```

**Fuji usage**

```fuji
let result = BeginMode3D(camera);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera3D` |

---

### EndMode3D

```c
RLAPI void EndMode3D()
```

**Fuji usage**

```fuji
let result = EndMode3D();
```

---

### BeginTextureMode

```c
RLAPI void BeginTextureMode(RenderTexture2D target)
```

**Fuji usage**

```fuji
let result = BeginTextureMode(target);
```

| Parameter | Type |
|-----------|------|
| `target` | `RenderTexture2D` |

---

### EndTextureMode

```c
RLAPI void EndTextureMode()
```

**Fuji usage**

```fuji
let result = EndTextureMode();
```

---

### BeginShaderMode

```c
RLAPI void BeginShaderMode(Shader shader)
```

**Fuji usage**

```fuji
let result = BeginShaderMode(shader);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |

---

### EndShaderMode

```c
RLAPI void EndShaderMode()
```

**Fuji usage**

```fuji
let result = EndShaderMode();
```

---

### BeginBlendMode

```c
RLAPI void BeginBlendMode(int mode)
```

**Fuji usage**

```fuji
let result = BeginBlendMode(mode);
```

| Parameter | Type |
|-----------|------|
| `mode` | `int` |

---

### EndBlendMode

```c
RLAPI void EndBlendMode()
```

**Fuji usage**

```fuji
let result = EndBlendMode();
```

---

### BeginScissorMode

```c
RLAPI void BeginScissorMode(int x, int y, int width, int height)
```

**Fuji usage**

```fuji
let result = BeginScissorMode(x, y, width, height);
```

| Parameter | Type |
|-----------|------|
| `x` | `int` |
| `y` | `int` |
| `width` | `int` |
| `height` | `int` |

---

### EndScissorMode

```c
RLAPI void EndScissorMode()
```

**Fuji usage**

```fuji
let result = EndScissorMode();
```

---

### BeginVrStereoMode

```c
RLAPI void BeginVrStereoMode(VrStereoConfig config)
```

**Fuji usage**

```fuji
let result = BeginVrStereoMode(config);
```

| Parameter | Type |
|-----------|------|
| `config` | `VrStereoConfig` |

---

### EndVrStereoMode

```c
RLAPI void EndVrStereoMode()
```

**Fuji usage**

```fuji
let result = EndVrStereoMode();
```

---

### LoadVrStereoConfig

```c
RLAPI VrStereoConfig LoadVrStereoConfig(VrDeviceInfo device)
```

**Fuji usage**

```fuji
let result = LoadVrStereoConfig(device);
```

| Parameter | Type |
|-----------|------|
| `device` | `VrDeviceInfo` |

---

### UnloadVrStereoConfig

```c
RLAPI void UnloadVrStereoConfig(VrStereoConfig config)
```

**Fuji usage**

```fuji
let result = UnloadVrStereoConfig(config);
```

| Parameter | Type |
|-----------|------|
| `config` | `VrStereoConfig` |

---

### LoadShader

```c
RLAPI Shader LoadShader(const char *vsFileName, const char *fsFileName)
```

**Fuji usage**

```fuji
let result = LoadShader(*vsFileName, *fsFileName);
```

| Parameter | Type |
|-----------|------|
| `*vsFileName` | `const char` |
| `*fsFileName` | `const char` |

---

### LoadShaderFromMemory

```c
RLAPI Shader LoadShaderFromMemory(const char *vsCode, const char *fsCode)
```

**Fuji usage**

```fuji
let result = LoadShaderFromMemory(*vsCode, *fsCode);
```

| Parameter | Type |
|-----------|------|
| `*vsCode` | `const char` |
| `*fsCode` | `const char` |

---

### IsShaderValid

```c
RLAPI bool IsShaderValid(Shader shader)
```

**Fuji usage**

```fuji
let result = IsShaderValid(shader);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |

---

### GetShaderLocation

```c
RLAPI int GetShaderLocation(Shader shader, const char *uniformName)
```

**Fuji usage**

```fuji
let result = GetShaderLocation(shader, *uniformName);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `*uniformName` | `const char` |

---

### GetShaderLocationAttrib

```c
RLAPI int GetShaderLocationAttrib(Shader shader, const char *attribName)
```

**Fuji usage**

```fuji
let result = GetShaderLocationAttrib(shader, *attribName);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `*attribName` | `const char` |

---

### SetShaderValue

```c
RLAPI void SetShaderValue(Shader shader, int locIndex, const void *value, int uniformType)
```

**Fuji usage**

```fuji
let result = SetShaderValue(shader, locIndex, *value, uniformType);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `locIndex` | `int` |
| `*value` | `const void` |
| `uniformType` | `int` |

---

### SetShaderValueV

```c
RLAPI void SetShaderValueV(Shader shader, int locIndex, const void *value, int uniformType, int count)
```

**Fuji usage**

```fuji
let result = SetShaderValueV(shader, locIndex, *value, uniformType, count);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `locIndex` | `int` |
| `*value` | `const void` |
| `uniformType` | `int` |
| `count` | `int` |

---

### SetShaderValueMatrix

```c
RLAPI void SetShaderValueMatrix(Shader shader, int locIndex, Matrix mat)
```

**Fuji usage**

```fuji
let result = SetShaderValueMatrix(shader, locIndex, mat);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `locIndex` | `int` |
| `mat` | `Matrix` |

---

### SetShaderValueTexture

```c
RLAPI void SetShaderValueTexture(Shader shader, int locIndex, Texture2D texture)
```

**Fuji usage**

```fuji
let result = SetShaderValueTexture(shader, locIndex, texture);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |
| `locIndex` | `int` |
| `texture` | `Texture2D` |

---

### UnloadShader

```c
RLAPI void UnloadShader(Shader shader)
```

**Fuji usage**

```fuji
let result = UnloadShader(shader);
```

| Parameter | Type |
|-----------|------|
| `shader` | `Shader` |

---

### GetScreenToWorldRay

```c
RLAPI Ray GetScreenToWorldRay(Vector2 position, Camera camera)
```

**Fuji usage**

```fuji
let result = GetScreenToWorldRay(position, camera);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `camera` | `Camera` |

---

### GetScreenToWorldRayEx

```c
RLAPI Ray GetScreenToWorldRayEx(Vector2 position, Camera camera, int width, int height)
```

**Fuji usage**

```fuji
let result = GetScreenToWorldRayEx(position, camera, width, height);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `camera` | `Camera` |
| `width` | `int` |
| `height` | `int` |

---

### GetWorldToScreen

```c
RLAPI Vector2 GetWorldToScreen(Vector3 position, Camera camera)
```

**Fuji usage**

```fuji
let result = GetWorldToScreen(position, camera);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `camera` | `Camera` |

---

### GetWorldToScreenEx

```c
RLAPI Vector2 GetWorldToScreenEx(Vector3 position, Camera camera, int width, int height)
```

**Fuji usage**

```fuji
let result = GetWorldToScreenEx(position, camera, width, height);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `camera` | `Camera` |
| `width` | `int` |
| `height` | `int` |

---

### GetWorldToScreen2D

```c
RLAPI Vector2 GetWorldToScreen2D(Vector2 position, Camera2D camera)
```

**Fuji usage**

```fuji
let result = GetWorldToScreen2D(position, camera);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `camera` | `Camera2D` |

---

### GetScreenToWorld2D

```c
RLAPI Vector2 GetScreenToWorld2D(Vector2 position, Camera2D camera)
```

**Fuji usage**

```fuji
let result = GetScreenToWorld2D(position, camera);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `camera` | `Camera2D` |

---

### GetCameraMatrix

```c
RLAPI Matrix GetCameraMatrix(Camera camera)
```

**Fuji usage**

```fuji
let result = GetCameraMatrix(camera);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera` |

---

### GetCameraMatrix2D

```c
RLAPI Matrix GetCameraMatrix2D(Camera2D camera)
```

**Fuji usage**

```fuji
let result = GetCameraMatrix2D(camera);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera2D` |

---

### SetTargetFPS

```c
RLAPI void SetTargetFPS(int fps)
```

**Fuji usage**

```fuji
let result = SetTargetFPS(fps);
```

| Parameter | Type |
|-----------|------|
| `fps` | `int` |

---

### GetFrameTime

```c
RLAPI float GetFrameTime()
```

**Fuji usage**

```fuji
let result = GetFrameTime();
```

---

### GetTime

```c
RLAPI double GetTime()
```

**Fuji usage**

```fuji
let result = GetTime();
```

---

### GetFPS

```c
RLAPI int GetFPS()
```

**Fuji usage**

```fuji
let result = GetFPS();
```

---

### SwapScreenBuffer

```c
RLAPI void SwapScreenBuffer()
```

**Fuji usage**

```fuji
let result = SwapScreenBuffer();
```

---

### PollInputEvents

```c
RLAPI void PollInputEvents()
```

**Fuji usage**

```fuji
let result = PollInputEvents();
```

---

### WaitTime

```c
RLAPI void WaitTime(double seconds)
```

**Fuji usage**

```fuji
let result = WaitTime(seconds);
```

| Parameter | Type |
|-----------|------|
| `seconds` | `double` |

---

### SetRandomSeed

```c
RLAPI void SetRandomSeed(unsigned int seed)
```

**Fuji usage**

```fuji
let result = SetRandomSeed(seed);
```

| Parameter | Type |
|-----------|------|
| `seed` | `unsigned int` |

---

### GetRandomValue

```c
RLAPI int GetRandomValue(int min, int max)
```

**Fuji usage**

```fuji
let result = GetRandomValue(min, max);
```

| Parameter | Type |
|-----------|------|
| `min` | `int` |
| `max` | `int` |

---

### UnloadRandomSequence

```c
RLAPI void UnloadRandomSequence(int *sequence)
```

**Fuji usage**

```fuji
let result = UnloadRandomSequence(*sequence);
```

| Parameter | Type |
|-----------|------|
| `*sequence` | `int` |

---

### TakeScreenshot

```c
RLAPI void TakeScreenshot(const char *fileName)
```

**Fuji usage**

```fuji
let result = TakeScreenshot(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### SetConfigFlags

```c
RLAPI void SetConfigFlags(unsigned int flags)
```

**Fuji usage**

```fuji
let result = SetConfigFlags(flags);
```

| Parameter | Type |
|-----------|------|
| `flags` | `unsigned int` |

---

### OpenURL

```c
RLAPI void OpenURL(const char *url)
```

**Fuji usage**

```fuji
let result = OpenURL(*url);
```

| Parameter | Type |
|-----------|------|
| `*url` | `const char` |

---

### SetTraceLogLevel

```c
RLAPI void SetTraceLogLevel(int logLevel)
```

**Fuji usage**

```fuji
let result = SetTraceLogLevel(logLevel);
```

| Parameter | Type |
|-----------|------|
| `logLevel` | `int` |

---

### TraceLog

```c
RLAPI void TraceLog(int logLevel, const char *text)
```

**Fuji usage**

```fuji
let result = TraceLog(logLevel, *text);
```

| Parameter | Type |
|-----------|------|
| `logLevel` | `int` |
| `*text` | `const char` |

---

### SetTraceLogCallback

```c
RLAPI void SetTraceLogCallback(TraceLogCallback callback)
```

**Fuji usage**

```fuji
let result = SetTraceLogCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `TraceLogCallback` |

---

### MemFree

```c
RLAPI void MemFree(void *ptr)
```

**Fuji usage**

```fuji
let result = MemFree(*ptr);
```

| Parameter | Type |
|-----------|------|
| `*ptr` | `void` |

---

### UnloadFileData

```c
RLAPI void UnloadFileData(unsigned char *data)
```

**Fuji usage**

```fuji
let result = UnloadFileData(*data);
```

| Parameter | Type |
|-----------|------|
| `*data` | `unsigned char` |

---

### SaveFileData

```c
RLAPI bool SaveFileData(const char *fileName, const void *data, int dataSize)
```

**Fuji usage**

```fuji
let result = SaveFileData(*fileName, *data, dataSize);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `*data` | `const void` |
| `dataSize` | `int` |

---

### ExportDataAsCode

```c
RLAPI bool ExportDataAsCode(const unsigned char *data, int dataSize, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportDataAsCode(*data, dataSize, *fileName);
```

| Parameter | Type |
|-----------|------|
| `*data` | `const unsigned char` |
| `dataSize` | `int` |
| `*fileName` | `const char` |

---

### UnloadFileText

```c
RLAPI void UnloadFileText(char *text)
```

**Fuji usage**

```fuji
let result = UnloadFileText(*text);
```

| Parameter | Type |
|-----------|------|
| `*text` | `char` |

---

### SaveFileText

```c
RLAPI bool SaveFileText(const char *fileName, const char *text)
```

**Fuji usage**

```fuji
let result = SaveFileText(*fileName, *text);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `*text` | `const char` |

---

### SetLoadFileDataCallback

```c
RLAPI void SetLoadFileDataCallback(LoadFileDataCallback callback)
```

**Fuji usage**

```fuji
let result = SetLoadFileDataCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `LoadFileDataCallback` |

---

### SetSaveFileDataCallback

```c
RLAPI void SetSaveFileDataCallback(SaveFileDataCallback callback)
```

**Fuji usage**

```fuji
let result = SetSaveFileDataCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `SaveFileDataCallback` |

---

### SetLoadFileTextCallback

```c
RLAPI void SetLoadFileTextCallback(LoadFileTextCallback callback)
```

**Fuji usage**

```fuji
let result = SetLoadFileTextCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `LoadFileTextCallback` |

---

### SetSaveFileTextCallback

```c
RLAPI void SetSaveFileTextCallback(SaveFileTextCallback callback)
```

**Fuji usage**

```fuji
let result = SetSaveFileTextCallback(callback);
```

| Parameter | Type |
|-----------|------|
| `callback` | `SaveFileTextCallback` |

---

### FileRename

```c
RLAPI int FileRename(const char *fileName, const char *fileRename)
```

**Fuji usage**

```fuji
let result = FileRename(*fileName, *fileRename);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `*fileRename` | `const char` |

---

### FileRemove

```c
RLAPI int FileRemove(const char *fileName)
```

**Fuji usage**

```fuji
let result = FileRemove(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### FileCopy

```c
RLAPI int FileCopy(const char *srcPath, const char *dstPath)
```

**Fuji usage**

```fuji
let result = FileCopy(*srcPath, *dstPath);
```

| Parameter | Type |
|-----------|------|
| `*srcPath` | `const char` |
| `*dstPath` | `const char` |

---

### FileMove

```c
RLAPI int FileMove(const char *srcPath, const char *dstPath)
```

**Fuji usage**

```fuji
let result = FileMove(*srcPath, *dstPath);
```

| Parameter | Type |
|-----------|------|
| `*srcPath` | `const char` |
| `*dstPath` | `const char` |

---

### FileTextReplace

```c
RLAPI int FileTextReplace(const char *fileName, const char *search, const char *replacement)
```

**Fuji usage**

```fuji
let result = FileTextReplace(*fileName, *search, *replacement);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `*search` | `const char` |
| `*replacement` | `const char` |

---

### FileTextFindIndex

```c
RLAPI int FileTextFindIndex(const char *fileName, const char *search)
```

**Fuji usage**

```fuji
let result = FileTextFindIndex(*fileName, *search);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `*search` | `const char` |

---

### FileExists

```c
RLAPI bool FileExists(const char *fileName)
```

**Fuji usage**

```fuji
let result = FileExists(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### DirectoryExists

```c
RLAPI bool DirectoryExists(const char *dirPath)
```

**Fuji usage**

```fuji
let result = DirectoryExists(*dirPath);
```

| Parameter | Type |
|-----------|------|
| `*dirPath` | `const char` |

---

### IsFileExtension

```c
RLAPI bool IsFileExtension(const char *fileName, const char *ext)
```

**Fuji usage**

```fuji
let result = IsFileExtension(*fileName, *ext);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `*ext` | `const char` |

---

### GetFileLength

```c
RLAPI int GetFileLength(const char *fileName)
```

**Fuji usage**

```fuji
let result = GetFileLength(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### GetFileModTime

```c
RLAPI long GetFileModTime(const char *fileName)
```

**Fuji usage**

```fuji
let result = GetFileModTime(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### MakeDirectory

```c
RLAPI int MakeDirectory(const char *dirPath)
```

**Fuji usage**

```fuji
let result = MakeDirectory(*dirPath);
```

| Parameter | Type |
|-----------|------|
| `*dirPath` | `const char` |

---

### ChangeDirectory

```c
RLAPI bool ChangeDirectory(const char *dirPath)
```

**Fuji usage**

```fuji
let result = ChangeDirectory(*dirPath);
```

| Parameter | Type |
|-----------|------|
| `*dirPath` | `const char` |

---

### IsPathFile

```c
RLAPI bool IsPathFile(const char *path)
```

**Fuji usage**

```fuji
let result = IsPathFile(*path);
```

| Parameter | Type |
|-----------|------|
| `*path` | `const char` |

---

### IsFileNameValid

```c
RLAPI bool IsFileNameValid(const char *fileName)
```

**Fuji usage**

```fuji
let result = IsFileNameValid(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### LoadDirectoryFiles

```c
RLAPI FilePathList LoadDirectoryFiles(const char *dirPath)
```

**Fuji usage**

```fuji
let result = LoadDirectoryFiles(*dirPath);
```

| Parameter | Type |
|-----------|------|
| `*dirPath` | `const char` |

---

### LoadDirectoryFilesEx

```c
RLAPI FilePathList LoadDirectoryFilesEx(const char *basePath, const char *filter, bool scanSubdirs)
```

**Fuji usage**

```fuji
let result = LoadDirectoryFilesEx(*basePath, *filter, scanSubdirs);
```

| Parameter | Type |
|-----------|------|
| `*basePath` | `const char` |
| `*filter` | `const char` |
| `scanSubdirs` | `bool` |

---

### UnloadDirectoryFiles

```c
RLAPI void UnloadDirectoryFiles(FilePathList files)
```

**Fuji usage**

```fuji
let result = UnloadDirectoryFiles(files);
```

| Parameter | Type |
|-----------|------|
| `files` | `FilePathList` |

---

### IsFileDropped

```c
RLAPI bool IsFileDropped()
```

**Fuji usage**

```fuji
let result = IsFileDropped();
```

---

### LoadDroppedFiles

```c
RLAPI FilePathList LoadDroppedFiles()
```

**Fuji usage**

```fuji
let result = LoadDroppedFiles();
```

---

### UnloadDroppedFiles

```c
RLAPI void UnloadDroppedFiles(FilePathList files)
```

**Fuji usage**

```fuji
let result = UnloadDroppedFiles(files);
```

| Parameter | Type |
|-----------|------|
| `files` | `FilePathList` |

---

### GetDirectoryFileCount

```c
RLAPI unsigned int GetDirectoryFileCount(const char *dirPath)
```

**Fuji usage**

```fuji
let result = GetDirectoryFileCount(*dirPath);
```

| Parameter | Type |
|-----------|------|
| `*dirPath` | `const char` |

---

### GetDirectoryFileCountEx

```c
RLAPI unsigned int GetDirectoryFileCountEx(const char *basePath, const char *filter, bool scanSubdirs)
```

**Fuji usage**

```fuji
let result = GetDirectoryFileCountEx(*basePath, *filter, scanSubdirs);
```

| Parameter | Type |
|-----------|------|
| `*basePath` | `const char` |
| `*filter` | `const char` |
| `scanSubdirs` | `bool` |

---

### ComputeCRC32

```c
RLAPI unsigned int ComputeCRC32(const unsigned char *data, int dataSize)
```

**Fuji usage**

```fuji
let result = ComputeCRC32(*data, dataSize);
```

| Parameter | Type |
|-----------|------|
| `*data` | `const unsigned char` |
| `dataSize` | `int` |

---

### LoadAutomationEventList

```c
RLAPI AutomationEventList LoadAutomationEventList(const char *fileName)
```

**Fuji usage**

```fuji
let result = LoadAutomationEventList(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### UnloadAutomationEventList

```c
RLAPI void UnloadAutomationEventList(AutomationEventList list)
```

**Fuji usage**

```fuji
let result = UnloadAutomationEventList(list);
```

| Parameter | Type |
|-----------|------|
| `list` | `AutomationEventList` |

---

### ExportAutomationEventList

```c
RLAPI bool ExportAutomationEventList(AutomationEventList list, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportAutomationEventList(list, *fileName);
```

| Parameter | Type |
|-----------|------|
| `list` | `AutomationEventList` |
| `*fileName` | `const char` |

---

### SetAutomationEventList

```c
RLAPI void SetAutomationEventList(AutomationEventList *list)
```

**Fuji usage**

```fuji
let result = SetAutomationEventList(*list);
```

| Parameter | Type |
|-----------|------|
| `*list` | `AutomationEventList` |

---

### SetAutomationEventBaseFrame

```c
RLAPI void SetAutomationEventBaseFrame(int frame)
```

**Fuji usage**

```fuji
let result = SetAutomationEventBaseFrame(frame);
```

| Parameter | Type |
|-----------|------|
| `frame` | `int` |

---

### StartAutomationEventRecording

```c
RLAPI void StartAutomationEventRecording()
```

**Fuji usage**

```fuji
let result = StartAutomationEventRecording();
```

---

### StopAutomationEventRecording

```c
RLAPI void StopAutomationEventRecording()
```

**Fuji usage**

```fuji
let result = StopAutomationEventRecording();
```

---

### PlayAutomationEvent

```c
RLAPI void PlayAutomationEvent(AutomationEvent event)
```

**Fuji usage**

```fuji
let result = PlayAutomationEvent(event);
```

| Parameter | Type |
|-----------|------|
| `event` | `AutomationEvent` |

---

### IsKeyPressed

```c
RLAPI bool IsKeyPressed(int key)
```

**Fuji usage**

```fuji
let result = IsKeyPressed(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsKeyPressedRepeat

```c
RLAPI bool IsKeyPressedRepeat(int key)
```

**Fuji usage**

```fuji
let result = IsKeyPressedRepeat(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsKeyDown

```c
RLAPI bool IsKeyDown(int key)
```

**Fuji usage**

```fuji
let result = IsKeyDown(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsKeyReleased

```c
RLAPI bool IsKeyReleased(int key)
```

**Fuji usage**

```fuji
let result = IsKeyReleased(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsKeyUp

```c
RLAPI bool IsKeyUp(int key)
```

**Fuji usage**

```fuji
let result = IsKeyUp(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### GetKeyPressed

```c
RLAPI int GetKeyPressed()
```

**Fuji usage**

```fuji
let result = GetKeyPressed();
```

---

### GetCharPressed

```c
RLAPI int GetCharPressed()
```

**Fuji usage**

```fuji
let result = GetCharPressed();
```

---

### SetExitKey

```c
RLAPI void SetExitKey(int key)
```

**Fuji usage**

```fuji
let result = SetExitKey(key);
```

| Parameter | Type |
|-----------|------|
| `key` | `int` |

---

### IsGamepadAvailable

```c
RLAPI bool IsGamepadAvailable(int gamepad)
```

**Fuji usage**

```fuji
let result = IsGamepadAvailable(gamepad);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |

---

### IsGamepadButtonPressed

```c
RLAPI bool IsGamepadButtonPressed(int gamepad, int button)
```

**Fuji usage**

```fuji
let result = IsGamepadButtonPressed(gamepad, button);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `button` | `int` |

---

### IsGamepadButtonDown

```c
RLAPI bool IsGamepadButtonDown(int gamepad, int button)
```

**Fuji usage**

```fuji
let result = IsGamepadButtonDown(gamepad, button);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `button` | `int` |

---

### IsGamepadButtonReleased

```c
RLAPI bool IsGamepadButtonReleased(int gamepad, int button)
```

**Fuji usage**

```fuji
let result = IsGamepadButtonReleased(gamepad, button);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `button` | `int` |

---

### IsGamepadButtonUp

```c
RLAPI bool IsGamepadButtonUp(int gamepad, int button)
```

**Fuji usage**

```fuji
let result = IsGamepadButtonUp(gamepad, button);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `button` | `int` |

---

### GetGamepadButtonPressed

```c
RLAPI int GetGamepadButtonPressed()
```

**Fuji usage**

```fuji
let result = GetGamepadButtonPressed();
```

---

### GetGamepadAxisCount

```c
RLAPI int GetGamepadAxisCount(int gamepad)
```

**Fuji usage**

```fuji
let result = GetGamepadAxisCount(gamepad);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |

---

### GetGamepadAxisMovement

```c
RLAPI float GetGamepadAxisMovement(int gamepad, int axis)
```

**Fuji usage**

```fuji
let result = GetGamepadAxisMovement(gamepad, axis);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `axis` | `int` |

---

### SetGamepadMappings

```c
RLAPI int SetGamepadMappings(const char *mappings)
```

**Fuji usage**

```fuji
let result = SetGamepadMappings(*mappings);
```

| Parameter | Type |
|-----------|------|
| `*mappings` | `const char` |

---

### SetGamepadVibration

```c
RLAPI void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
```

**Fuji usage**

```fuji
let result = SetGamepadVibration(gamepad, leftMotor, rightMotor, duration);
```

| Parameter | Type |
|-----------|------|
| `gamepad` | `int` |
| `leftMotor` | `float` |
| `rightMotor` | `float` |
| `duration` | `float` |

---

### IsMouseButtonPressed

```c
RLAPI bool IsMouseButtonPressed(int button)
```

**Fuji usage**

```fuji
let result = IsMouseButtonPressed(button);
```

| Parameter | Type |
|-----------|------|
| `button` | `int` |

---

### IsMouseButtonDown

```c
RLAPI bool IsMouseButtonDown(int button)
```

**Fuji usage**

```fuji
let result = IsMouseButtonDown(button);
```

| Parameter | Type |
|-----------|------|
| `button` | `int` |

---

### IsMouseButtonReleased

```c
RLAPI bool IsMouseButtonReleased(int button)
```

**Fuji usage**

```fuji
let result = IsMouseButtonReleased(button);
```

| Parameter | Type |
|-----------|------|
| `button` | `int` |

---

### IsMouseButtonUp

```c
RLAPI bool IsMouseButtonUp(int button)
```

**Fuji usage**

```fuji
let result = IsMouseButtonUp(button);
```

| Parameter | Type |
|-----------|------|
| `button` | `int` |

---

### GetMouseX

```c
RLAPI int GetMouseX()
```

**Fuji usage**

```fuji
let result = GetMouseX();
```

---

### GetMouseY

```c
RLAPI int GetMouseY()
```

**Fuji usage**

```fuji
let result = GetMouseY();
```

---

### GetMousePosition

```c
RLAPI Vector2 GetMousePosition()
```

**Fuji usage**

```fuji
let result = GetMousePosition();
```

---

### GetMouseDelta

```c
RLAPI Vector2 GetMouseDelta()
```

**Fuji usage**

```fuji
let result = GetMouseDelta();
```

---

### SetMousePosition

```c
RLAPI void SetMousePosition(int x, int y)
```

**Fuji usage**

```fuji
let result = SetMousePosition(x, y);
```

| Parameter | Type |
|-----------|------|
| `x` | `int` |
| `y` | `int` |

---

### SetMouseOffset

```c
RLAPI void SetMouseOffset(int offsetX, int offsetY)
```

**Fuji usage**

```fuji
let result = SetMouseOffset(offsetX, offsetY);
```

| Parameter | Type |
|-----------|------|
| `offsetX` | `int` |
| `offsetY` | `int` |

---

### SetMouseScale

```c
RLAPI void SetMouseScale(float scaleX, float scaleY)
```

**Fuji usage**

```fuji
let result = SetMouseScale(scaleX, scaleY);
```

| Parameter | Type |
|-----------|------|
| `scaleX` | `float` |
| `scaleY` | `float` |

---

### GetMouseWheelMove

```c
RLAPI float GetMouseWheelMove()
```

**Fuji usage**

```fuji
let result = GetMouseWheelMove();
```

---

### GetMouseWheelMoveV

```c
RLAPI Vector2 GetMouseWheelMoveV()
```

**Fuji usage**

```fuji
let result = GetMouseWheelMoveV();
```

---

### SetMouseCursor

```c
RLAPI void SetMouseCursor(int cursor)
```

**Fuji usage**

```fuji
let result = SetMouseCursor(cursor);
```

| Parameter | Type |
|-----------|------|
| `cursor` | `int` |

---

### GetTouchX

```c
RLAPI int GetTouchX()
```

**Fuji usage**

```fuji
let result = GetTouchX();
```

---

### GetTouchY

```c
RLAPI int GetTouchY()
```

**Fuji usage**

```fuji
let result = GetTouchY();
```

---

### GetTouchPosition

```c
RLAPI Vector2 GetTouchPosition(int index)
```

**Fuji usage**

```fuji
let result = GetTouchPosition(index);
```

| Parameter | Type |
|-----------|------|
| `index` | `int` |

---

### GetTouchPointId

```c
RLAPI int GetTouchPointId(int index)
```

**Fuji usage**

```fuji
let result = GetTouchPointId(index);
```

| Parameter | Type |
|-----------|------|
| `index` | `int` |

---

### GetTouchPointCount

```c
RLAPI int GetTouchPointCount()
```

**Fuji usage**

```fuji
let result = GetTouchPointCount();
```

---

### SetGesturesEnabled

```c
RLAPI void SetGesturesEnabled(unsigned int flags)
```

**Fuji usage**

```fuji
let result = SetGesturesEnabled(flags);
```

| Parameter | Type |
|-----------|------|
| `flags` | `unsigned int` |

---

### IsGestureDetected

```c
RLAPI bool IsGestureDetected(unsigned int gesture)
```

**Fuji usage**

```fuji
let result = IsGestureDetected(gesture);
```

| Parameter | Type |
|-----------|------|
| `gesture` | `unsigned int` |

---

### GetGestureDetected

```c
RLAPI int GetGestureDetected()
```

**Fuji usage**

```fuji
let result = GetGestureDetected();
```

---

### GetGestureHoldDuration

```c
RLAPI float GetGestureHoldDuration()
```

**Fuji usage**

```fuji
let result = GetGestureHoldDuration();
```

---

### GetGestureDragVector

```c
RLAPI Vector2 GetGestureDragVector()
```

**Fuji usage**

```fuji
let result = GetGestureDragVector();
```

---

### GetGestureDragAngle

```c
RLAPI float GetGestureDragAngle()
```

**Fuji usage**

```fuji
let result = GetGestureDragAngle();
```

---

### GetGesturePinchVector

```c
RLAPI Vector2 GetGesturePinchVector()
```

**Fuji usage**

```fuji
let result = GetGesturePinchVector();
```

---

### GetGesturePinchAngle

```c
RLAPI float GetGesturePinchAngle()
```

**Fuji usage**

```fuji
let result = GetGesturePinchAngle();
```

---

### UpdateCamera

```c
RLAPI void UpdateCamera(Camera *camera, int mode)
```

**Fuji usage**

```fuji
let result = UpdateCamera(*camera, mode);
```

| Parameter | Type |
|-----------|------|
| `*camera` | `Camera` |
| `mode` | `int` |

---

### UpdateCameraPro

```c
RLAPI void UpdateCameraPro(Camera *camera, Vector3 movement, Vector3 rotation, float zoom)
```

**Fuji usage**

```fuji
let result = UpdateCameraPro(*camera, movement, rotation, zoom);
```

| Parameter | Type |
|-----------|------|
| `*camera` | `Camera` |
| `movement` | `Vector3` |
| `rotation` | `Vector3` |
| `zoom` | `float` |

---

### SetShapesTexture

```c
RLAPI void SetShapesTexture(Texture2D texture, Rectangle source)
```

**Fuji usage**

```fuji
let result = SetShapesTexture(texture, source);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `source` | `Rectangle` |

---

### GetShapesTexture

```c
RLAPI Texture2D GetShapesTexture()
```

**Fuji usage**

```fuji
let result = GetShapesTexture();
```

---

### GetShapesTextureRectangle

```c
RLAPI Rectangle GetShapesTextureRectangle()
```

**Fuji usage**

```fuji
let result = GetShapesTextureRectangle();
```

---

### DrawPixel

```c
RLAPI void DrawPixel(int posX, int posY, Color color)
```

**Fuji usage**

```fuji
let result = DrawPixel(posX, posY, color);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `color` | `Color` |

---

### DrawPixelV

```c
RLAPI void DrawPixelV(Vector2 position, Color color)
```

**Fuji usage**

```fuji
let result = DrawPixelV(position, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `color` | `Color` |

---

### DrawLine

```c
RLAPI void DrawLine(int startPosX, int startPosY, int endPosX, int endPosY, Color color)
```

**Fuji usage**

```fuji
let result = DrawLine(startPosX, startPosY, endPosX, endPosY, color);
```

| Parameter | Type |
|-----------|------|
| `startPosX` | `int` |
| `startPosY` | `int` |
| `endPosX` | `int` |
| `endPosY` | `int` |
| `color` | `Color` |

---

### DrawLineV

```c
RLAPI void DrawLineV(Vector2 startPos, Vector2 endPos, Color color)
```

**Fuji usage**

```fuji
let result = DrawLineV(startPos, endPos, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `color` | `Color` |

---

### DrawLineEx

```c
RLAPI void DrawLineEx(Vector2 startPos, Vector2 endPos, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawLineEx(startPos, endPos, thick, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawLineStrip

```c
RLAPI void DrawLineStrip(const Vector2 *points, int pointCount, Color color)
```

**Fuji usage**

```fuji
let result = DrawLineStrip(*points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### DrawLineBezier

```c
RLAPI void DrawLineBezier(Vector2 startPos, Vector2 endPos, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawLineBezier(startPos, endPos, thick, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawLineDashed

```c
RLAPI void DrawLineDashed(Vector2 startPos, Vector2 endPos, int dashSize, int spaceSize, Color color)
```

**Fuji usage**

```fuji
let result = DrawLineDashed(startPos, endPos, dashSize, spaceSize, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `dashSize` | `int` |
| `spaceSize` | `int` |
| `color` | `Color` |

---

### DrawCircle

```c
RLAPI void DrawCircle(int centerX, int centerY, float radius, Color color)
```

**Fuji usage**

```fuji
let result = DrawCircle(centerX, centerY, radius, color);
```

| Parameter | Type |
|-----------|------|
| `centerX` | `int` |
| `centerY` | `int` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawCircleV

```c
RLAPI void DrawCircleV(Vector2 center, float radius, Color color)
```

**Fuji usage**

```fuji
let result = DrawCircleV(center, radius, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawCircleGradient

```c
RLAPI void DrawCircleGradient(Vector2 center, float radius, Color inner, Color outer)
```

**Fuji usage**

```fuji
let result = DrawCircleGradient(center, radius, inner, outer);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `inner` | `Color` |
| `outer` | `Color` |

---

### DrawCircleSector

```c
RLAPI void DrawCircleSector(Vector2 center, float radius, float startAngle, float endAngle, int segments, Color color)
```

**Fuji usage**

```fuji
let result = DrawCircleSector(center, radius, startAngle, endAngle, segments, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `startAngle` | `float` |
| `endAngle` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawCircleSectorLines

```c
RLAPI void DrawCircleSectorLines(Vector2 center, float radius, float startAngle, float endAngle, int segments, Color color)
```

**Fuji usage**

```fuji
let result = DrawCircleSectorLines(center, radius, startAngle, endAngle, segments, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `startAngle` | `float` |
| `endAngle` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawCircleLines

```c
RLAPI void DrawCircleLines(int centerX, int centerY, float radius, Color color)
```

**Fuji usage**

```fuji
let result = DrawCircleLines(centerX, centerY, radius, color);
```

| Parameter | Type |
|-----------|------|
| `centerX` | `int` |
| `centerY` | `int` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawCircleLinesV

```c
RLAPI void DrawCircleLinesV(Vector2 center, float radius, Color color)
```

**Fuji usage**

```fuji
let result = DrawCircleLinesV(center, radius, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawEllipse

```c
RLAPI void DrawEllipse(int centerX, int centerY, float radiusH, float radiusV, Color color)
```

**Fuji usage**

```fuji
let result = DrawEllipse(centerX, centerY, radiusH, radiusV, color);
```

| Parameter | Type |
|-----------|------|
| `centerX` | `int` |
| `centerY` | `int` |
| `radiusH` | `float` |
| `radiusV` | `float` |
| `color` | `Color` |

---

### DrawEllipseV

```c
RLAPI void DrawEllipseV(Vector2 center, float radiusH, float radiusV, Color color)
```

**Fuji usage**

```fuji
let result = DrawEllipseV(center, radiusH, radiusV, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radiusH` | `float` |
| `radiusV` | `float` |
| `color` | `Color` |

---

### DrawEllipseLines

```c
RLAPI void DrawEllipseLines(int centerX, int centerY, float radiusH, float radiusV, Color color)
```

**Fuji usage**

```fuji
let result = DrawEllipseLines(centerX, centerY, radiusH, radiusV, color);
```

| Parameter | Type |
|-----------|------|
| `centerX` | `int` |
| `centerY` | `int` |
| `radiusH` | `float` |
| `radiusV` | `float` |
| `color` | `Color` |

---

### DrawEllipseLinesV

```c
RLAPI void DrawEllipseLinesV(Vector2 center, float radiusH, float radiusV, Color color)
```

**Fuji usage**

```fuji
let result = DrawEllipseLinesV(center, radiusH, radiusV, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radiusH` | `float` |
| `radiusV` | `float` |
| `color` | `Color` |

---

### DrawRing

```c
RLAPI void DrawRing(Vector2 center, float innerRadius, float outerRadius, float startAngle, float endAngle, int segments, Color color)
```

**Fuji usage**

```fuji
let result = DrawRing(center, innerRadius, outerRadius, startAngle, endAngle, segments, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `innerRadius` | `float` |
| `outerRadius` | `float` |
| `startAngle` | `float` |
| `endAngle` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawRingLines

```c
RLAPI void DrawRingLines(Vector2 center, float innerRadius, float outerRadius, float startAngle, float endAngle, int segments, Color color)
```

**Fuji usage**

```fuji
let result = DrawRingLines(center, innerRadius, outerRadius, startAngle, endAngle, segments, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `innerRadius` | `float` |
| `outerRadius` | `float` |
| `startAngle` | `float` |
| `endAngle` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawRectangle

```c
RLAPI void DrawRectangle(int posX, int posY, int width, int height, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectangle(posX, posY, width, height, color);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### DrawRectangleV

```c
RLAPI void DrawRectangleV(Vector2 position, Vector2 size, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectangleV(position, size, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector2` |
| `size` | `Vector2` |
| `color` | `Color` |

---

### DrawRectangleRec

```c
RLAPI void DrawRectangleRec(Rectangle rec, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectangleRec(rec, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `color` | `Color` |

---

### DrawRectanglePro

```c
RLAPI void DrawRectanglePro(Rectangle rec, Vector2 origin, float rotation, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectanglePro(rec, origin, rotation, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `color` | `Color` |

---

### DrawRectangleGradientV

```c
RLAPI void DrawRectangleGradientV(int posX, int posY, int width, int height, Color top, Color bottom)
```

**Fuji usage**

```fuji
let result = DrawRectangleGradientV(posX, posY, width, height, top, bottom);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `top` | `Color` |
| `bottom` | `Color` |

---

### DrawRectangleGradientH

```c
RLAPI void DrawRectangleGradientH(int posX, int posY, int width, int height, Color left, Color right)
```

**Fuji usage**

```fuji
let result = DrawRectangleGradientH(posX, posY, width, height, left, right);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `left` | `Color` |
| `right` | `Color` |

---

### DrawRectangleGradientEx

```c
RLAPI void DrawRectangleGradientEx(Rectangle rec, Color topLeft, Color bottomLeft, Color bottomRight, Color topRight)
```

**Fuji usage**

```fuji
let result = DrawRectangleGradientEx(rec, topLeft, bottomLeft, bottomRight, topRight);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `topLeft` | `Color` |
| `bottomLeft` | `Color` |
| `bottomRight` | `Color` |
| `topRight` | `Color` |

---

### DrawRectangleLines

```c
RLAPI void DrawRectangleLines(int posX, int posY, int width, int height, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectangleLines(posX, posY, width, height, color);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### DrawRectangleLinesEx

```c
RLAPI void DrawRectangleLinesEx(Rectangle rec, float lineThick, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectangleLinesEx(rec, lineThick, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `lineThick` | `float` |
| `color` | `Color` |

---

### DrawRectangleRounded

```c
RLAPI void DrawRectangleRounded(Rectangle rec, float roundness, int segments, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectangleRounded(rec, roundness, segments, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `roundness` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawRectangleRoundedLines

```c
RLAPI void DrawRectangleRoundedLines(Rectangle rec, float roundness, int segments, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectangleRoundedLines(rec, roundness, segments, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `roundness` | `float` |
| `segments` | `int` |
| `color` | `Color` |

---

### DrawRectangleRoundedLinesEx

```c
RLAPI void DrawRectangleRoundedLinesEx(Rectangle rec, float roundness, int segments, float lineThick, Color color)
```

**Fuji usage**

```fuji
let result = DrawRectangleRoundedLinesEx(rec, roundness, segments, lineThick, color);
```

| Parameter | Type |
|-----------|------|
| `rec` | `Rectangle` |
| `roundness` | `float` |
| `segments` | `int` |
| `lineThick` | `float` |
| `color` | `Color` |

---

### DrawTriangle

```c
RLAPI void DrawTriangle(Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Fuji usage**

```fuji
let result = DrawTriangle(v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `color` | `Color` |

---

### DrawTriangleGradient

```c
RLAPI void DrawTriangleGradient(Vector2 v1, Vector2 v2, Vector2 v3, Color c1, Color c2, Color c3)
```

**Fuji usage**

```fuji
let result = DrawTriangleGradient(v1, v2, v3, c1, c2, c3);
```

| Parameter | Type |
|-----------|------|
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `c1` | `Color` |
| `c2` | `Color` |
| `c3` | `Color` |

---

### DrawTriangleLines

```c
RLAPI void DrawTriangleLines(Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Fuji usage**

```fuji
let result = DrawTriangleLines(v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `color` | `Color` |

---

### DrawTriangleFan

```c
RLAPI void DrawTriangleFan(const Vector2 *points, int pointCount, Color color)
```

**Fuji usage**

```fuji
let result = DrawTriangleFan(*points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### DrawTriangleStrip

```c
RLAPI void DrawTriangleStrip(const Vector2 *points, int pointCount, Color color)
```

**Fuji usage**

```fuji
let result = DrawTriangleStrip(*points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### DrawPoly

```c
RLAPI void DrawPoly(Vector2 center, int sides, float radius, float rotation, Color color)
```

**Fuji usage**

```fuji
let result = DrawPoly(center, sides, radius, rotation, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `sides` | `int` |
| `radius` | `float` |
| `rotation` | `float` |
| `color` | `Color` |

---

### DrawPolyLines

```c
RLAPI void DrawPolyLines(Vector2 center, int sides, float radius, float rotation, Color color)
```

**Fuji usage**

```fuji
let result = DrawPolyLines(center, sides, radius, rotation, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `sides` | `int` |
| `radius` | `float` |
| `rotation` | `float` |
| `color` | `Color` |

---

### DrawPolyLinesEx

```c
RLAPI void DrawPolyLinesEx(Vector2 center, int sides, float radius, float rotation, float lineThick, Color color)
```

**Fuji usage**

```fuji
let result = DrawPolyLinesEx(center, sides, radius, rotation, lineThick, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `sides` | `int` |
| `radius` | `float` |
| `rotation` | `float` |
| `lineThick` | `float` |
| `color` | `Color` |

---

### DrawSplineLinear

```c
RLAPI void DrawSplineLinear(const Vector2 *points, int pointCount, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineLinear(*points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineBasis

```c
RLAPI void DrawSplineBasis(const Vector2 *points, int pointCount, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineBasis(*points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineCatmullRom

```c
RLAPI void DrawSplineCatmullRom(const Vector2 *points, int pointCount, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineCatmullRom(*points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineBezierQuadratic

```c
RLAPI void DrawSplineBezierQuadratic(const Vector2 *points, int pointCount, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineBezierQuadratic(*points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineBezierCubic

```c
RLAPI void DrawSplineBezierCubic(const Vector2 *points, int pointCount, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineBezierCubic(*points, pointCount, thick, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentLinear

```c
RLAPI void DrawSplineSegmentLinear(Vector2 p1, Vector2 p2, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineSegmentLinear(p1, p2, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentBasis

```c
RLAPI void DrawSplineSegmentBasis(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineSegmentBasis(p1, p2, p3, p4, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |
| `p4` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentCatmullRom

```c
RLAPI void DrawSplineSegmentCatmullRom(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineSegmentCatmullRom(p1, p2, p3, p4, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |
| `p4` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentBezierQuadratic

```c
RLAPI void DrawSplineSegmentBezierQuadratic(Vector2 p1, Vector2 c2, Vector2 p3, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineSegmentBezierQuadratic(p1, c2, p3, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `c2` | `Vector2` |
| `p3` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### DrawSplineSegmentBezierCubic

```c
RLAPI void DrawSplineSegmentBezierCubic(Vector2 p1, Vector2 c2, Vector2 c3, Vector2 p4, float thick, Color color)
```

**Fuji usage**

```fuji
let result = DrawSplineSegmentBezierCubic(p1, c2, c3, p4, thick, color);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `c2` | `Vector2` |
| `c3` | `Vector2` |
| `p4` | `Vector2` |
| `thick` | `float` |
| `color` | `Color` |

---

### GetSplinePointLinear

```c
RLAPI Vector2 GetSplinePointLinear(Vector2 startPos, Vector2 endPos, float t)
```

**Fuji usage**

```fuji
let result = GetSplinePointLinear(startPos, endPos, t);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector2` |
| `endPos` | `Vector2` |
| `t` | `float` |

---

### GetSplinePointBasis

```c
RLAPI Vector2 GetSplinePointBasis(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float t)
```

**Fuji usage**

```fuji
let result = GetSplinePointBasis(p1, p2, p3, p4, t);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |
| `p4` | `Vector2` |
| `t` | `float` |

---

### GetSplinePointCatmullRom

```c
RLAPI Vector2 GetSplinePointCatmullRom(Vector2 p1, Vector2 p2, Vector2 p3, Vector2 p4, float t)
```

**Fuji usage**

```fuji
let result = GetSplinePointCatmullRom(p1, p2, p3, p4, t);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |
| `p4` | `Vector2` |
| `t` | `float` |

---

### GetSplinePointBezierQuadratic

```c
RLAPI Vector2 GetSplinePointBezierQuadratic(Vector2 p1, Vector2 c2, Vector2 p3, float t)
```

**Fuji usage**

```fuji
let result = GetSplinePointBezierQuadratic(p1, c2, p3, t);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `c2` | `Vector2` |
| `p3` | `Vector2` |
| `t` | `float` |

---

### GetSplinePointBezierCubic

```c
RLAPI Vector2 GetSplinePointBezierCubic(Vector2 p1, Vector2 c2, Vector2 c3, Vector2 p4, float t)
```

**Fuji usage**

```fuji
let result = GetSplinePointBezierCubic(p1, c2, c3, p4, t);
```

| Parameter | Type |
|-----------|------|
| `p1` | `Vector2` |
| `c2` | `Vector2` |
| `c3` | `Vector2` |
| `p4` | `Vector2` |
| `t` | `float` |

---

### CheckCollisionRecs

```c
RLAPI bool CheckCollisionRecs(Rectangle rec1, Rectangle rec2)
```

**Fuji usage**

```fuji
let result = CheckCollisionRecs(rec1, rec2);
```

| Parameter | Type |
|-----------|------|
| `rec1` | `Rectangle` |
| `rec2` | `Rectangle` |

---

### CheckCollisionCircles

```c
RLAPI bool CheckCollisionCircles(Vector2 center1, float radius1, Vector2 center2, float radius2)
```

**Fuji usage**

```fuji
let result = CheckCollisionCircles(center1, radius1, center2, radius2);
```

| Parameter | Type |
|-----------|------|
| `center1` | `Vector2` |
| `radius1` | `float` |
| `center2` | `Vector2` |
| `radius2` | `float` |

---

### CheckCollisionCircleRec

```c
RLAPI bool CheckCollisionCircleRec(Vector2 center, float radius, Rectangle rec)
```

**Fuji usage**

```fuji
let result = CheckCollisionCircleRec(center, radius, rec);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `rec` | `Rectangle` |

---

### CheckCollisionCircleLine

```c
RLAPI bool CheckCollisionCircleLine(Vector2 center, float radius, Vector2 p1, Vector2 p2)
```

**Fuji usage**

```fuji
let result = CheckCollisionCircleLine(center, radius, p1, p2);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector2` |
| `radius` | `float` |
| `p1` | `Vector2` |
| `p2` | `Vector2` |

---

### CheckCollisionPointRec

```c
RLAPI bool CheckCollisionPointRec(Vector2 point, Rectangle rec)
```

**Fuji usage**

```fuji
let result = CheckCollisionPointRec(point, rec);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `rec` | `Rectangle` |

---

### CheckCollisionPointCircle

```c
RLAPI bool CheckCollisionPointCircle(Vector2 point, Vector2 center, float radius)
```

**Fuji usage**

```fuji
let result = CheckCollisionPointCircle(point, center, radius);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `center` | `Vector2` |
| `radius` | `float` |

---

### CheckCollisionPointTriangle

```c
RLAPI bool CheckCollisionPointTriangle(Vector2 point, Vector2 p1, Vector2 p2, Vector2 p3)
```

**Fuji usage**

```fuji
let result = CheckCollisionPointTriangle(point, p1, p2, p3);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `p3` | `Vector2` |

---

### CheckCollisionPointLine

```c
RLAPI bool CheckCollisionPointLine(Vector2 point, Vector2 p1, Vector2 p2, int threshold)
```

**Fuji usage**

```fuji
let result = CheckCollisionPointLine(point, p1, p2, threshold);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `p1` | `Vector2` |
| `p2` | `Vector2` |
| `threshold` | `int` |

---

### CheckCollisionPointPoly

```c
RLAPI bool CheckCollisionPointPoly(Vector2 point, const Vector2 *points, int pointCount)
```

**Fuji usage**

```fuji
let result = CheckCollisionPointPoly(point, *points, pointCount);
```

| Parameter | Type |
|-----------|------|
| `point` | `Vector2` |
| `*points` | `const Vector2` |
| `pointCount` | `int` |

---

### CheckCollisionLines

```c
RLAPI bool CheckCollisionLines(Vector2 startPos1, Vector2 endPos1, Vector2 startPos2, Vector2 endPos2, Vector2 *collisionPoint)
```

**Fuji usage**

```fuji
let result = CheckCollisionLines(startPos1, endPos1, startPos2, endPos2, *collisionPoint);
```

| Parameter | Type |
|-----------|------|
| `startPos1` | `Vector2` |
| `endPos1` | `Vector2` |
| `startPos2` | `Vector2` |
| `endPos2` | `Vector2` |
| `*collisionPoint` | `Vector2` |

---

### GetCollisionRec

```c
RLAPI Rectangle GetCollisionRec(Rectangle rec1, Rectangle rec2)
```

**Fuji usage**

```fuji
let result = GetCollisionRec(rec1, rec2);
```

| Parameter | Type |
|-----------|------|
| `rec1` | `Rectangle` |
| `rec2` | `Rectangle` |

---

### LoadImage

```c
RLAPI Image LoadImage(const char *fileName)
```

**Fuji usage**

```fuji
let result = LoadImage(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### LoadImageRaw

```c
RLAPI Image LoadImageRaw(const char *fileName, int width, int height, int format, int headerSize)
```

**Fuji usage**

```fuji
let result = LoadImageRaw(*fileName, width, height, format, headerSize);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `width` | `int` |
| `height` | `int` |
| `format` | `int` |
| `headerSize` | `int` |

---

### LoadImageAnim

```c
RLAPI Image LoadImageAnim(const char *fileName, int *frames)
```

**Fuji usage**

```fuji
let result = LoadImageAnim(*fileName, *frames);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `*frames` | `int` |

---

### LoadImageAnimFromMemory

```c
RLAPI Image LoadImageAnimFromMemory(const char *fileType, const unsigned char *fileData, int dataSize, int *frames)
```

**Fuji usage**

```fuji
let result = LoadImageAnimFromMemory(*fileType, *fileData, dataSize, *frames);
```

| Parameter | Type |
|-----------|------|
| `*fileType` | `const char` |
| `*fileData` | `const unsigned char` |
| `dataSize` | `int` |
| `*frames` | `int` |

---

### LoadImageFromMemory

```c
RLAPI Image LoadImageFromMemory(const char *fileType, const unsigned char *fileData, int dataSize)
```

**Fuji usage**

```fuji
let result = LoadImageFromMemory(*fileType, *fileData, dataSize);
```

| Parameter | Type |
|-----------|------|
| `*fileType` | `const char` |
| `*fileData` | `const unsigned char` |
| `dataSize` | `int` |

---

### LoadImageFromTexture

```c
RLAPI Image LoadImageFromTexture(Texture2D texture)
```

**Fuji usage**

```fuji
let result = LoadImageFromTexture(texture);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |

---

### LoadImageFromScreen

```c
RLAPI Image LoadImageFromScreen()
```

**Fuji usage**

```fuji
let result = LoadImageFromScreen();
```

---

### IsImageValid

```c
RLAPI bool IsImageValid(Image image)
```

**Fuji usage**

```fuji
let result = IsImageValid(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### UnloadImage

```c
RLAPI void UnloadImage(Image image)
```

**Fuji usage**

```fuji
let result = UnloadImage(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### ExportImage

```c
RLAPI bool ExportImage(Image image, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportImage(image, *fileName);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `*fileName` | `const char` |

---

### ExportImageAsCode

```c
RLAPI bool ExportImageAsCode(Image image, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportImageAsCode(image, *fileName);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `*fileName` | `const char` |

---

### GenImageColor

```c
RLAPI Image GenImageColor(int width, int height, Color color)
```

**Fuji usage**

```fuji
let result = GenImageColor(width, height, color);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### GenImageGradientLinear

```c
RLAPI Image GenImageGradientLinear(int width, int height, int direction, Color start, Color end)
```

**Fuji usage**

```fuji
let result = GenImageGradientLinear(width, height, direction, start, end);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `direction` | `int` |
| `start` | `Color` |
| `end` | `Color` |

---

### GenImageGradientRadial

```c
RLAPI Image GenImageGradientRadial(int width, int height, float density, Color inner, Color outer)
```

**Fuji usage**

```fuji
let result = GenImageGradientRadial(width, height, density, inner, outer);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `density` | `float` |
| `inner` | `Color` |
| `outer` | `Color` |

---

### GenImageGradientSquare

```c
RLAPI Image GenImageGradientSquare(int width, int height, float density, Color inner, Color outer)
```

**Fuji usage**

```fuji
let result = GenImageGradientSquare(width, height, density, inner, outer);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `density` | `float` |
| `inner` | `Color` |
| `outer` | `Color` |

---

### GenImageChecked

```c
RLAPI Image GenImageChecked(int width, int height, int checksX, int checksY, Color col1, Color col2)
```

**Fuji usage**

```fuji
let result = GenImageChecked(width, height, checksX, checksY, col1, col2);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `checksX` | `int` |
| `checksY` | `int` |
| `col1` | `Color` |
| `col2` | `Color` |

---

### GenImageWhiteNoise

```c
RLAPI Image GenImageWhiteNoise(int width, int height, float factor)
```

**Fuji usage**

```fuji
let result = GenImageWhiteNoise(width, height, factor);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `factor` | `float` |

---

### GenImagePerlinNoise

```c
RLAPI Image GenImagePerlinNoise(int width, int height, int offsetX, int offsetY, float scale)
```

**Fuji usage**

```fuji
let result = GenImagePerlinNoise(width, height, offsetX, offsetY, scale);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `offsetX` | `int` |
| `offsetY` | `int` |
| `scale` | `float` |

---

### GenImageCellular

```c
RLAPI Image GenImageCellular(int width, int height, int tileSize)
```

**Fuji usage**

```fuji
let result = GenImageCellular(width, height, tileSize);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `tileSize` | `int` |

---

### GenImageText

```c
RLAPI Image GenImageText(int width, int height, const char *text)
```

**Fuji usage**

```fuji
let result = GenImageText(width, height, *text);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `*text` | `const char` |

---

### ImageCopy

```c
RLAPI Image ImageCopy(Image image)
```

**Fuji usage**

```fuji
let result = ImageCopy(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### ImageFromImage

```c
RLAPI Image ImageFromImage(Image image, Rectangle rec)
```

**Fuji usage**

```fuji
let result = ImageFromImage(image, rec);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `rec` | `Rectangle` |

---

### ImageFromChannel

```c
RLAPI Image ImageFromChannel(Image image, int selectedChannel)
```

**Fuji usage**

```fuji
let result = ImageFromChannel(image, selectedChannel);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `selectedChannel` | `int` |

---

### ImageText

```c
RLAPI Image ImageText(const char *text, int fontSize, Color color)
```

**Fuji usage**

```fuji
let result = ImageText(*text, fontSize, color);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |
| `fontSize` | `int` |
| `color` | `Color` |

---

### ImageTextEx

```c
RLAPI Image ImageTextEx(Font font, const char *text, float fontSize, float spacing, Color tint)
```

**Fuji usage**

```fuji
let result = ImageTextEx(font, *text, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `*text` | `const char` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### ImageFormat

```c
RLAPI void ImageFormat(Image *image, int newFormat)
```

**Fuji usage**

```fuji
let result = ImageFormat(*image, newFormat);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `newFormat` | `int` |

---

### ImageToPOT

```c
RLAPI void ImageToPOT(Image *image, Color fill)
```

**Fuji usage**

```fuji
let result = ImageToPOT(*image, fill);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `fill` | `Color` |

---

### ImageCrop

```c
RLAPI void ImageCrop(Image *image, Rectangle crop)
```

**Fuji usage**

```fuji
let result = ImageCrop(*image, crop);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `crop` | `Rectangle` |

---

### ImageAlphaCrop

```c
RLAPI void ImageAlphaCrop(Image *image, float threshold)
```

**Fuji usage**

```fuji
let result = ImageAlphaCrop(*image, threshold);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `threshold` | `float` |

---

### ImageAlphaClear

```c
RLAPI void ImageAlphaClear(Image *image, Color color, float threshold)
```

**Fuji usage**

```fuji
let result = ImageAlphaClear(*image, color, threshold);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `color` | `Color` |
| `threshold` | `float` |

---

### ImageAlphaMask

```c
RLAPI void ImageAlphaMask(Image *image, Image alphaMask)
```

**Fuji usage**

```fuji
let result = ImageAlphaMask(*image, alphaMask);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `alphaMask` | `Image` |

---

### ImageAlphaPremultiply

```c
RLAPI void ImageAlphaPremultiply(Image *image)
```

**Fuji usage**

```fuji
let result = ImageAlphaPremultiply(*image);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |

---

### ImageBlurGaussian

```c
RLAPI void ImageBlurGaussian(Image *image, int blurSize)
```

**Fuji usage**

```fuji
let result = ImageBlurGaussian(*image, blurSize);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `blurSize` | `int` |

---

### ImageKernelConvolution

```c
RLAPI void ImageKernelConvolution(Image *image, const float *kernel, int kernelSize)
```

**Fuji usage**

```fuji
let result = ImageKernelConvolution(*image, *kernel, kernelSize);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `*kernel` | `const float` |
| `kernelSize` | `int` |

---

### ImageResize

```c
RLAPI void ImageResize(Image *image, int newWidth, int newHeight)
```

**Fuji usage**

```fuji
let result = ImageResize(*image, newWidth, newHeight);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `newWidth` | `int` |
| `newHeight` | `int` |

---

### ImageResizeNN

```c
RLAPI void ImageResizeNN(Image *image, int newWidth, int newHeight)
```

**Fuji usage**

```fuji
let result = ImageResizeNN(*image, newWidth, newHeight);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `newWidth` | `int` |
| `newHeight` | `int` |

---

### ImageResizeCanvas

```c
RLAPI void ImageResizeCanvas(Image *image, int newWidth, int newHeight, int offsetX, int offsetY, Color fill)
```

**Fuji usage**

```fuji
let result = ImageResizeCanvas(*image, newWidth, newHeight, offsetX, offsetY, fill);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `newWidth` | `int` |
| `newHeight` | `int` |
| `offsetX` | `int` |
| `offsetY` | `int` |
| `fill` | `Color` |

---

### ImageMipmaps

```c
RLAPI void ImageMipmaps(Image *image)
```

**Fuji usage**

```fuji
let result = ImageMipmaps(*image);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |

---

### ImageDither

```c
RLAPI void ImageDither(Image *image, int rBpp, int gBpp, int bBpp, int aBpp)
```

**Fuji usage**

```fuji
let result = ImageDither(*image, rBpp, gBpp, bBpp, aBpp);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `rBpp` | `int` |
| `gBpp` | `int` |
| `bBpp` | `int` |
| `aBpp` | `int` |

---

### ImageFlipVertical

```c
RLAPI void ImageFlipVertical(Image *image)
```

**Fuji usage**

```fuji
let result = ImageFlipVertical(*image);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |

---

### ImageFlipHorizontal

```c
RLAPI void ImageFlipHorizontal(Image *image)
```

**Fuji usage**

```fuji
let result = ImageFlipHorizontal(*image);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |

---

### ImageRotate

```c
RLAPI void ImageRotate(Image *image, int degrees)
```

**Fuji usage**

```fuji
let result = ImageRotate(*image, degrees);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `degrees` | `int` |

---

### ImageRotateCW

```c
RLAPI void ImageRotateCW(Image *image)
```

**Fuji usage**

```fuji
let result = ImageRotateCW(*image);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |

---

### ImageRotateCCW

```c
RLAPI void ImageRotateCCW(Image *image)
```

**Fuji usage**

```fuji
let result = ImageRotateCCW(*image);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |

---

### ImageColorTint

```c
RLAPI void ImageColorTint(Image *image, Color color)
```

**Fuji usage**

```fuji
let result = ImageColorTint(*image, color);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `color` | `Color` |

---

### ImageColorInvert

```c
RLAPI void ImageColorInvert(Image *image)
```

**Fuji usage**

```fuji
let result = ImageColorInvert(*image);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |

---

### ImageColorGrayscale

```c
RLAPI void ImageColorGrayscale(Image *image)
```

**Fuji usage**

```fuji
let result = ImageColorGrayscale(*image);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |

---

### ImageColorContrast

```c
RLAPI void ImageColorContrast(Image *image, int contrast)
```

**Fuji usage**

```fuji
let result = ImageColorContrast(*image, contrast);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `contrast` | `int` |

---

### ImageColorBrightness

```c
RLAPI void ImageColorBrightness(Image *image, int brightness)
```

**Fuji usage**

```fuji
let result = ImageColorBrightness(*image, brightness);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `brightness` | `int` |

---

### ImageColorReplace

```c
RLAPI void ImageColorReplace(Image *image, Color color, Color replace)
```

**Fuji usage**

```fuji
let result = ImageColorReplace(*image, color, replace);
```

| Parameter | Type |
|-----------|------|
| `*image` | `Image` |
| `color` | `Color` |
| `replace` | `Color` |

---

### UnloadImageColors

```c
RLAPI void UnloadImageColors(Color *colors)
```

**Fuji usage**

```fuji
let result = UnloadImageColors(*colors);
```

| Parameter | Type |
|-----------|------|
| `*colors` | `Color` |

---

### UnloadImagePalette

```c
RLAPI void UnloadImagePalette(Color *colors)
```

**Fuji usage**

```fuji
let result = UnloadImagePalette(*colors);
```

| Parameter | Type |
|-----------|------|
| `*colors` | `Color` |

---

### GetImageAlphaBorder

```c
RLAPI Rectangle GetImageAlphaBorder(Image image, float threshold)
```

**Fuji usage**

```fuji
let result = GetImageAlphaBorder(image, threshold);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `threshold` | `float` |

---

### GetImageColor

```c
RLAPI Color GetImageColor(Image image, int x, int y)
```

**Fuji usage**

```fuji
let result = GetImageColor(image, x, y);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `x` | `int` |
| `y` | `int` |

---

### ImageClearBackground

```c
RLAPI void ImageClearBackground(Image *dst, Color color)
```

**Fuji usage**

```fuji
let result = ImageClearBackground(*dst, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `color` | `Color` |

---

### ImageDrawPixel

```c
RLAPI void ImageDrawPixel(Image *dst, int posX, int posY, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawPixel(*dst, posX, posY, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `posX` | `int` |
| `posY` | `int` |
| `color` | `Color` |

---

### ImageDrawPixelV

```c
RLAPI void ImageDrawPixelV(Image *dst, Vector2 position, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawPixelV(*dst, position, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `position` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawLine

```c
RLAPI void ImageDrawLine(Image *dst, int startPosX, int startPosY, int endPosX, int endPosY, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawLine(*dst, startPosX, startPosY, endPosX, endPosY, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `startPosX` | `int` |
| `startPosY` | `int` |
| `endPosX` | `int` |
| `endPosY` | `int` |
| `color` | `Color` |

---

### ImageDrawLineV

```c
RLAPI void ImageDrawLineV(Image *dst, Vector2 start, Vector2 end, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawLineV(*dst, start, end, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `start` | `Vector2` |
| `end` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawLineEx

```c
RLAPI void ImageDrawLineEx(Image *dst, Vector2 start, Vector2 end, int thick, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawLineEx(*dst, start, end, thick, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `start` | `Vector2` |
| `end` | `Vector2` |
| `thick` | `int` |
| `color` | `Color` |

---

### ImageDrawCircle

```c
RLAPI void ImageDrawCircle(Image *dst, int centerX, int centerY, int radius, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawCircle(*dst, centerX, centerY, radius, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `centerX` | `int` |
| `centerY` | `int` |
| `radius` | `int` |
| `color` | `Color` |

---

### ImageDrawCircleV

```c
RLAPI void ImageDrawCircleV(Image *dst, Vector2 center, int radius, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawCircleV(*dst, center, radius, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `center` | `Vector2` |
| `radius` | `int` |
| `color` | `Color` |

---

### ImageDrawCircleLines

```c
RLAPI void ImageDrawCircleLines(Image *dst, int centerX, int centerY, int radius, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawCircleLines(*dst, centerX, centerY, radius, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `centerX` | `int` |
| `centerY` | `int` |
| `radius` | `int` |
| `color` | `Color` |

---

### ImageDrawCircleLinesV

```c
RLAPI void ImageDrawCircleLinesV(Image *dst, Vector2 center, int radius, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawCircleLinesV(*dst, center, radius, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `center` | `Vector2` |
| `radius` | `int` |
| `color` | `Color` |

---

### ImageDrawRectangle

```c
RLAPI void ImageDrawRectangle(Image *dst, int posX, int posY, int width, int height, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawRectangle(*dst, posX, posY, width, height, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### ImageDrawRectangleV

```c
RLAPI void ImageDrawRectangleV(Image *dst, Vector2 position, Vector2 size, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawRectangleV(*dst, position, size, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `position` | `Vector2` |
| `size` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawRectangleRec

```c
RLAPI void ImageDrawRectangleRec(Image *dst, Rectangle rec, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawRectangleRec(*dst, rec, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `rec` | `Rectangle` |
| `color` | `Color` |

---

### ImageDrawRectangleLines

```c
RLAPI void ImageDrawRectangleLines(Image *dst, int posX, int posY, int width, int height, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawRectangleLines(*dst, posX, posY, width, height, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `posX` | `int` |
| `posY` | `int` |
| `width` | `int` |
| `height` | `int` |
| `color` | `Color` |

---

### ImageDrawRectangleLinesEx

```c
RLAPI void ImageDrawRectangleLinesEx(Image *dst, Rectangle rec, int thick, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawRectangleLinesEx(*dst, rec, thick, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `rec` | `Rectangle` |
| `thick` | `int` |
| `color` | `Color` |

---

### ImageDrawTriangle

```c
RLAPI void ImageDrawTriangle(Image *dst, Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawTriangle(*dst, v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawTriangleGradient

```c
RLAPI void ImageDrawTriangleGradient(Image *dst, Vector2 v1, Vector2 v2, Vector2 v3, Color c1, Color c2, Color c3)
```

**Fuji usage**

```fuji
let result = ImageDrawTriangleGradient(*dst, v1, v2, v3, c1, c2, c3);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `c1` | `Color` |
| `c2` | `Color` |
| `c3` | `Color` |

---

### ImageDrawTriangleLines

```c
RLAPI void ImageDrawTriangleLines(Image *dst, Vector2 v1, Vector2 v2, Vector2 v3, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawTriangleLines(*dst, v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `v1` | `Vector2` |
| `v2` | `Vector2` |
| `v3` | `Vector2` |
| `color` | `Color` |

---

### ImageDrawTriangleFan

```c
RLAPI void ImageDrawTriangleFan(Image *dst, const Vector2 *points, int pointCount, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawTriangleFan(*dst, *points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### ImageDrawTriangleStrip

```c
RLAPI void ImageDrawTriangleStrip(Image *dst, const Vector2 *points, int pointCount, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawTriangleStrip(*dst, *points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `*points` | `const Vector2` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### ImageDraw

```c
RLAPI void ImageDraw(Image *dst, Image src, Rectangle srcRec, Rectangle dstRec, Color tint)
```

**Fuji usage**

```fuji
let result = ImageDraw(*dst, src, srcRec, dstRec, tint);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `src` | `Image` |
| `srcRec` | `Rectangle` |
| `dstRec` | `Rectangle` |
| `tint` | `Color` |

---

### ImageDrawText

```c
RLAPI void ImageDrawText(Image *dst, const char *text, int posX, int posY, int fontSize, Color color)
```

**Fuji usage**

```fuji
let result = ImageDrawText(*dst, *text, posX, posY, fontSize, color);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `*text` | `const char` |
| `posX` | `int` |
| `posY` | `int` |
| `fontSize` | `int` |
| `color` | `Color` |

---

### ImageDrawTextEx

```c
RLAPI void ImageDrawTextEx(Image *dst, Font font, const char *text, Vector2 position, float fontSize, float spacing, Color tint)
```

**Fuji usage**

```fuji
let result = ImageDrawTextEx(*dst, font, *text, position, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `Image` |
| `font` | `Font` |
| `*text` | `const char` |
| `position` | `Vector2` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### LoadTexture

```c
RLAPI Texture2D LoadTexture(const char *fileName)
```

**Fuji usage**

```fuji
let result = LoadTexture(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### LoadTextureFromImage

```c
RLAPI Texture2D LoadTextureFromImage(Image image)
```

**Fuji usage**

```fuji
let result = LoadTextureFromImage(image);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |

---

### LoadTextureCubemap

```c
RLAPI TextureCubemap LoadTextureCubemap(Image image, int layout)
```

**Fuji usage**

```fuji
let result = LoadTextureCubemap(image, layout);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `layout` | `int` |

---

### LoadRenderTexture

```c
RLAPI RenderTexture2D LoadRenderTexture(int width, int height)
```

**Fuji usage**

```fuji
let result = LoadRenderTexture(width, height);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |

---

### IsTextureValid

```c
RLAPI bool IsTextureValid(Texture2D texture)
```

**Fuji usage**

```fuji
let result = IsTextureValid(texture);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |

---

### UnloadTexture

```c
RLAPI void UnloadTexture(Texture2D texture)
```

**Fuji usage**

```fuji
let result = UnloadTexture(texture);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |

---

### IsRenderTextureValid

```c
RLAPI bool IsRenderTextureValid(RenderTexture2D target)
```

**Fuji usage**

```fuji
let result = IsRenderTextureValid(target);
```

| Parameter | Type |
|-----------|------|
| `target` | `RenderTexture2D` |

---

### UnloadRenderTexture

```c
RLAPI void UnloadRenderTexture(RenderTexture2D target)
```

**Fuji usage**

```fuji
let result = UnloadRenderTexture(target);
```

| Parameter | Type |
|-----------|------|
| `target` | `RenderTexture2D` |

---

### UpdateTexture

```c
RLAPI void UpdateTexture(Texture2D texture, const void *pixels)
```

**Fuji usage**

```fuji
let result = UpdateTexture(texture, *pixels);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `*pixels` | `const void` |

---

### UpdateTextureRec

```c
RLAPI void UpdateTextureRec(Texture2D texture, Rectangle rec, const void *pixels)
```

**Fuji usage**

```fuji
let result = UpdateTextureRec(texture, rec, *pixels);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `rec` | `Rectangle` |
| `*pixels` | `const void` |

---

### GenTextureMipmaps

```c
RLAPI void GenTextureMipmaps(Texture2D *texture)
```

**Fuji usage**

```fuji
let result = GenTextureMipmaps(*texture);
```

| Parameter | Type |
|-----------|------|
| `*texture` | `Texture2D` |

---

### SetTextureFilter

```c
RLAPI void SetTextureFilter(Texture2D texture, int filter)
```

**Fuji usage**

```fuji
let result = SetTextureFilter(texture, filter);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `filter` | `int` |

---

### SetTextureWrap

```c
RLAPI void SetTextureWrap(Texture2D texture, int wrap)
```

**Fuji usage**

```fuji
let result = SetTextureWrap(texture, wrap);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `wrap` | `int` |

---

### DrawTexture

```c
RLAPI void DrawTexture(Texture2D texture, int posX, int posY, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTexture(texture, posX, posY, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `posX` | `int` |
| `posY` | `int` |
| `tint` | `Color` |

---

### DrawTextureV

```c
RLAPI void DrawTextureV(Texture2D texture, Vector2 position, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTextureV(texture, position, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `position` | `Vector2` |
| `tint` | `Color` |

---

### DrawTextureEx

```c
RLAPI void DrawTextureEx(Texture2D texture, Vector2 position, float rotation, float scale, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTextureEx(texture, position, rotation, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `position` | `Vector2` |
| `rotation` | `float` |
| `scale` | `float` |
| `tint` | `Color` |

---

### DrawTextureRec

```c
RLAPI void DrawTextureRec(Texture2D texture, Rectangle source, Vector2 position, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTextureRec(texture, source, position, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `source` | `Rectangle` |
| `position` | `Vector2` |
| `tint` | `Color` |

---

### DrawTexturePro

```c
RLAPI void DrawTexturePro(Texture2D texture, Rectangle source, Rectangle dest, Vector2 origin, float rotation, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTexturePro(texture, source, dest, origin, rotation, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `source` | `Rectangle` |
| `dest` | `Rectangle` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `tint` | `Color` |

---

### DrawTextureNPatch

```c
RLAPI void DrawTextureNPatch(Texture2D texture, NPatchInfo nPatchInfo, Rectangle dest, Vector2 origin, float rotation, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTextureNPatch(texture, nPatchInfo, dest, origin, rotation, tint);
```

| Parameter | Type |
|-----------|------|
| `texture` | `Texture2D` |
| `nPatchInfo` | `NPatchInfo` |
| `dest` | `Rectangle` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `tint` | `Color` |

---

### ColorIsEqual

```c
RLAPI bool ColorIsEqual(Color col1, Color col2)
```

**Fuji usage**

```fuji
let result = ColorIsEqual(col1, col2);
```

| Parameter | Type |
|-----------|------|
| `col1` | `Color` |
| `col2` | `Color` |

---

### Fade

```c
RLAPI Color Fade(Color color, float alpha)
```

**Fuji usage**

```fuji
let result = Fade(color, alpha);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `alpha` | `float` |

---

### ColorToInt

```c
RLAPI int ColorToInt(Color color)
```

**Fuji usage**

```fuji
let result = ColorToInt(color);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |

---

### ColorNormalize

```c
RLAPI Vector4 ColorNormalize(Color color)
```

**Fuji usage**

```fuji
let result = ColorNormalize(color);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |

---

### ColorFromNormalized

```c
RLAPI Color ColorFromNormalized(Vector4 normalized)
```

**Fuji usage**

```fuji
let result = ColorFromNormalized(normalized);
```

| Parameter | Type |
|-----------|------|
| `normalized` | `Vector4` |

---

### ColorToHSV

```c
RLAPI Vector3 ColorToHSV(Color color)
```

**Fuji usage**

```fuji
let result = ColorToHSV(color);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |

---

### ColorFromHSV

```c
RLAPI Color ColorFromHSV(float hue, float saturation, float value)
```

**Fuji usage**

```fuji
let result = ColorFromHSV(hue, saturation, value);
```

| Parameter | Type |
|-----------|------|
| `hue` | `float` |
| `saturation` | `float` |
| `value` | `float` |

---

### ColorTint

```c
RLAPI Color ColorTint(Color color, Color tint)
```

**Fuji usage**

```fuji
let result = ColorTint(color, tint);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `tint` | `Color` |

---

### ColorBrightness

```c
RLAPI Color ColorBrightness(Color color, float factor)
```

**Fuji usage**

```fuji
let result = ColorBrightness(color, factor);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `factor` | `float` |

---

### ColorContrast

```c
RLAPI Color ColorContrast(Color color, float contrast)
```

**Fuji usage**

```fuji
let result = ColorContrast(color, contrast);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `contrast` | `float` |

---

### ColorAlpha

```c
RLAPI Color ColorAlpha(Color color, float alpha)
```

**Fuji usage**

```fuji
let result = ColorAlpha(color, alpha);
```

| Parameter | Type |
|-----------|------|
| `color` | `Color` |
| `alpha` | `float` |

---

### ColorAlphaBlend

```c
RLAPI Color ColorAlphaBlend(Color dst, Color src, Color tint)
```

**Fuji usage**

```fuji
let result = ColorAlphaBlend(dst, src, tint);
```

| Parameter | Type |
|-----------|------|
| `dst` | `Color` |
| `src` | `Color` |
| `tint` | `Color` |

---

### ColorLerp

```c
RLAPI Color ColorLerp(Color color1, Color color2, float factor)
```

**Fuji usage**

```fuji
let result = ColorLerp(color1, color2, factor);
```

| Parameter | Type |
|-----------|------|
| `color1` | `Color` |
| `color2` | `Color` |
| `factor` | `float` |

---

### GetColor

```c
RLAPI Color GetColor(unsigned int hexValue)
```

**Fuji usage**

```fuji
let result = GetColor(hexValue);
```

| Parameter | Type |
|-----------|------|
| `hexValue` | `unsigned int` |

---

### GetPixelColor

```c
RLAPI Color GetPixelColor(void *srcPtr, int format)
```

**Fuji usage**

```fuji
let result = GetPixelColor(*srcPtr, format);
```

| Parameter | Type |
|-----------|------|
| `*srcPtr` | `void` |
| `format` | `int` |

---

### SetPixelColor

```c
RLAPI void SetPixelColor(void *dstPtr, Color color, int format)
```

**Fuji usage**

```fuji
let result = SetPixelColor(*dstPtr, color, format);
```

| Parameter | Type |
|-----------|------|
| `*dstPtr` | `void` |
| `color` | `Color` |
| `format` | `int` |

---

### GetPixelDataSize

```c
RLAPI int GetPixelDataSize(int width, int height, int format)
```

**Fuji usage**

```fuji
let result = GetPixelDataSize(width, height, format);
```

| Parameter | Type |
|-----------|------|
| `width` | `int` |
| `height` | `int` |
| `format` | `int` |

---

### GetFontDefault

```c
RLAPI Font GetFontDefault()
```

**Fuji usage**

```fuji
let result = GetFontDefault();
```

---

### LoadFont

```c
RLAPI Font LoadFont(const char *fileName)
```

**Fuji usage**

```fuji
let result = LoadFont(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### LoadFontEx

```c
RLAPI Font LoadFontEx(const char *fileName, int fontSize, const int *codepoints, int codepointCount)
```

**Fuji usage**

```fuji
let result = LoadFontEx(*fileName, fontSize, *codepoints, codepointCount);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |
| `fontSize` | `int` |
| `*codepoints` | `const int` |
| `codepointCount` | `int` |

---

### LoadFontFromImage

```c
RLAPI Font LoadFontFromImage(Image image, Color key, int firstChar)
```

**Fuji usage**

```fuji
let result = LoadFontFromImage(image, key, firstChar);
```

| Parameter | Type |
|-----------|------|
| `image` | `Image` |
| `key` | `Color` |
| `firstChar` | `int` |

---

### LoadFontFromMemory

```c
RLAPI Font LoadFontFromMemory(const char *fileType, const unsigned char *fileData, int dataSize, int fontSize, const int *codepoints, int codepointCount)
```

**Fuji usage**

```fuji
let result = LoadFontFromMemory(*fileType, *fileData, dataSize, fontSize, *codepoints, codepointCount);
```

| Parameter | Type |
|-----------|------|
| `*fileType` | `const char` |
| `*fileData` | `const unsigned char` |
| `dataSize` | `int` |
| `fontSize` | `int` |
| `*codepoints` | `const int` |
| `codepointCount` | `int` |

---

### IsFontValid

```c
RLAPI bool IsFontValid(Font font)
```

**Fuji usage**

```fuji
let result = IsFontValid(font);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |

---

### GenImageFontAtlas

```c
RLAPI Image GenImageFontAtlas(const GlyphInfo *glyphs, Rectangle **glyphRecs, int glyphCount, int fontSize, int padding, int packMethod)
```

**Fuji usage**

```fuji
let result = GenImageFontAtlas(*glyphs, **glyphRecs, glyphCount, fontSize, padding, packMethod);
```

| Parameter | Type |
|-----------|------|
| `*glyphs` | `const GlyphInfo` |
| `**glyphRecs` | `Rectangle` |
| `glyphCount` | `int` |
| `fontSize` | `int` |
| `padding` | `int` |
| `packMethod` | `int` |

---

### UnloadFontData

```c
RLAPI void UnloadFontData(GlyphInfo *glyphs, int glyphCount)
```

**Fuji usage**

```fuji
let result = UnloadFontData(*glyphs, glyphCount);
```

| Parameter | Type |
|-----------|------|
| `*glyphs` | `GlyphInfo` |
| `glyphCount` | `int` |

---

### UnloadFont

```c
RLAPI void UnloadFont(Font font)
```

**Fuji usage**

```fuji
let result = UnloadFont(font);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |

---

### ExportFontAsCode

```c
RLAPI bool ExportFontAsCode(Font font, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportFontAsCode(font, *fileName);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `*fileName` | `const char` |

---

### DrawFPS

```c
RLAPI void DrawFPS(int posX, int posY)
```

**Fuji usage**

```fuji
let result = DrawFPS(posX, posY);
```

| Parameter | Type |
|-----------|------|
| `posX` | `int` |
| `posY` | `int` |

---

### DrawText

```c
RLAPI void DrawText(const char *text, int posX, int posY, int fontSize, Color color)
```

**Fuji usage**

```fuji
let result = DrawText(*text, posX, posY, fontSize, color);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |
| `posX` | `int` |
| `posY` | `int` |
| `fontSize` | `int` |
| `color` | `Color` |

---

### DrawTextEx

```c
RLAPI void DrawTextEx(Font font, const char *text, Vector2 position, float fontSize, float spacing, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTextEx(font, *text, position, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `*text` | `const char` |
| `position` | `Vector2` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### DrawTextPro

```c
RLAPI void DrawTextPro(Font font, const char *text, Vector2 position, Vector2 origin, float rotation, float fontSize, float spacing, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTextPro(font, *text, position, origin, rotation, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `*text` | `const char` |
| `position` | `Vector2` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### DrawTextCodepoint

```c
RLAPI void DrawTextCodepoint(Font font, int codepoint, Vector2 position, float fontSize, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTextCodepoint(font, codepoint, position, fontSize, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoint` | `int` |
| `position` | `Vector2` |
| `fontSize` | `float` |
| `tint` | `Color` |

---

### DrawTextCodepoints

```c
RLAPI void DrawTextCodepoints(Font font, const int *codepoints, int codepointCount, Vector2 position, float fontSize, float spacing, Color tint)
```

**Fuji usage**

```fuji
let result = DrawTextCodepoints(font, *codepoints, codepointCount, position, fontSize, spacing, tint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `*codepoints` | `const int` |
| `codepointCount` | `int` |
| `position` | `Vector2` |
| `fontSize` | `float` |
| `spacing` | `float` |
| `tint` | `Color` |

---

### SetTextLineSpacing

```c
RLAPI void SetTextLineSpacing(int spacing)
```

**Fuji usage**

```fuji
let result = SetTextLineSpacing(spacing);
```

| Parameter | Type |
|-----------|------|
| `spacing` | `int` |

---

### MeasureText

```c
RLAPI int MeasureText(const char *text, int fontSize)
```

**Fuji usage**

```fuji
let result = MeasureText(*text, fontSize);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |
| `fontSize` | `int` |

---

### MeasureTextEx

```c
RLAPI Vector2 MeasureTextEx(Font font, const char *text, float fontSize, float spacing)
```

**Fuji usage**

```fuji
let result = MeasureTextEx(font, *text, fontSize, spacing);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `*text` | `const char` |
| `fontSize` | `float` |
| `spacing` | `float` |

---

### MeasureTextCodepoints

```c
RLAPI Vector2 MeasureTextCodepoints(Font font, const int *codepoints, int length, float fontSize, float spacing)
```

**Fuji usage**

```fuji
let result = MeasureTextCodepoints(font, *codepoints, length, fontSize, spacing);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `*codepoints` | `const int` |
| `length` | `int` |
| `fontSize` | `float` |
| `spacing` | `float` |

---

### GetGlyphIndex

```c
RLAPI int GetGlyphIndex(Font font, int codepoint)
```

**Fuji usage**

```fuji
let result = GetGlyphIndex(font, codepoint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoint` | `int` |

---

### GetGlyphInfo

```c
RLAPI GlyphInfo GetGlyphInfo(Font font, int codepoint)
```

**Fuji usage**

```fuji
let result = GetGlyphInfo(font, codepoint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoint` | `int` |

---

### GetGlyphAtlasRec

```c
RLAPI Rectangle GetGlyphAtlasRec(Font font, int codepoint)
```

**Fuji usage**

```fuji
let result = GetGlyphAtlasRec(font, codepoint);
```

| Parameter | Type |
|-----------|------|
| `font` | `Font` |
| `codepoint` | `int` |

---

### UnloadUTF8

```c
RLAPI void UnloadUTF8(char *text)
```

**Fuji usage**

```fuji
let result = UnloadUTF8(*text);
```

| Parameter | Type |
|-----------|------|
| `*text` | `char` |

---

### UnloadCodepoints

```c
RLAPI void UnloadCodepoints(int *codepoints)
```

**Fuji usage**

```fuji
let result = UnloadCodepoints(*codepoints);
```

| Parameter | Type |
|-----------|------|
| `*codepoints` | `int` |

---

### GetCodepointCount

```c
RLAPI int GetCodepointCount(const char *text)
```

**Fuji usage**

```fuji
let result = GetCodepointCount(*text);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |

---

### GetCodepoint

```c
RLAPI int GetCodepoint(const char *text, int *codepointSize)
```

**Fuji usage**

```fuji
let result = GetCodepoint(*text, *codepointSize);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |
| `*codepointSize` | `int` |

---

### GetCodepointNext

```c
RLAPI int GetCodepointNext(const char *text, int *codepointSize)
```

**Fuji usage**

```fuji
let result = GetCodepointNext(*text, *codepointSize);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |
| `*codepointSize` | `int` |

---

### GetCodepointPrevious

```c
RLAPI int GetCodepointPrevious(const char *text, int *codepointSize)
```

**Fuji usage**

```fuji
let result = GetCodepointPrevious(*text, *codepointSize);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |
| `*codepointSize` | `int` |

---

### UnloadTextLines

```c
RLAPI void UnloadTextLines(char **text, int lineCount)
```

**Fuji usage**

```fuji
let result = UnloadTextLines(**text, lineCount);
```

| Parameter | Type |
|-----------|------|
| `**text` | `char` |
| `lineCount` | `int` |

---

### TextCopy

```c
RLAPI int TextCopy(char *dst, const char *src)
```

**Fuji usage**

```fuji
let result = TextCopy(*dst, *src);
```

| Parameter | Type |
|-----------|------|
| `*dst` | `char` |
| `*src` | `const char` |

---

### TextIsEqual

```c
RLAPI bool TextIsEqual(const char *text1, const char *text2)
```

**Fuji usage**

```fuji
let result = TextIsEqual(*text1, *text2);
```

| Parameter | Type |
|-----------|------|
| `*text1` | `const char` |
| `*text2` | `const char` |

---

### TextLength

```c
RLAPI unsigned int TextLength(const char *text)
```

**Fuji usage**

```fuji
let result = TextLength(*text);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |

---

### TextAppend

```c
RLAPI void TextAppend(char *text, const char *append, int *position)
```

**Fuji usage**

```fuji
let result = TextAppend(*text, *append, *position);
```

| Parameter | Type |
|-----------|------|
| `*text` | `char` |
| `*append` | `const char` |
| `*position` | `int` |

---

### TextFindIndex

```c
RLAPI int TextFindIndex(const char *text, const char *search)
```

**Fuji usage**

```fuji
let result = TextFindIndex(*text, *search);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |
| `*search` | `const char` |

---

### TextToInteger

```c
RLAPI int TextToInteger(const char *text)
```

**Fuji usage**

```fuji
let result = TextToInteger(*text);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |

---

### TextToFloat

```c
RLAPI float TextToFloat(const char *text)
```

**Fuji usage**

```fuji
let result = TextToFloat(*text);
```

| Parameter | Type |
|-----------|------|
| `*text` | `const char` |

---

### DrawLine3D

```c
RLAPI void DrawLine3D(Vector3 startPos, Vector3 endPos, Color color)
```

**Fuji usage**

```fuji
let result = DrawLine3D(startPos, endPos, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `color` | `Color` |

---

### DrawPoint3D

```c
RLAPI void DrawPoint3D(Vector3 position, Color color)
```

**Fuji usage**

```fuji
let result = DrawPoint3D(position, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `color` | `Color` |

---

### DrawCircle3D

```c
RLAPI void DrawCircle3D(Vector3 center, float radius, Vector3 rotationAxis, float rotationAngle, Color color)
```

**Fuji usage**

```fuji
let result = DrawCircle3D(center, radius, rotationAxis, rotationAngle, color);
```

| Parameter | Type |
|-----------|------|
| `center` | `Vector3` |
| `radius` | `float` |
| `rotationAxis` | `Vector3` |
| `rotationAngle` | `float` |
| `color` | `Color` |

---

### DrawTriangle3D

```c
RLAPI void DrawTriangle3D(Vector3 v1, Vector3 v2, Vector3 v3, Color color)
```

**Fuji usage**

```fuji
let result = DrawTriangle3D(v1, v2, v3, color);
```

| Parameter | Type |
|-----------|------|
| `v1` | `Vector3` |
| `v2` | `Vector3` |
| `v3` | `Vector3` |
| `color` | `Color` |

---

### DrawTriangleStrip3D

```c
RLAPI void DrawTriangleStrip3D(const Vector3 *points, int pointCount, Color color)
```

**Fuji usage**

```fuji
let result = DrawTriangleStrip3D(*points, pointCount, color);
```

| Parameter | Type |
|-----------|------|
| `*points` | `const Vector3` |
| `pointCount` | `int` |
| `color` | `Color` |

---

### DrawCube

```c
RLAPI void DrawCube(Vector3 position, float width, float height, float length, Color color)
```

**Fuji usage**

```fuji
let result = DrawCube(position, width, height, length, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `width` | `float` |
| `height` | `float` |
| `length` | `float` |
| `color` | `Color` |

---

### DrawCubeV

```c
RLAPI void DrawCubeV(Vector3 position, Vector3 size, Color color)
```

**Fuji usage**

```fuji
let result = DrawCubeV(position, size, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `size` | `Vector3` |
| `color` | `Color` |

---

### DrawCubeWires

```c
RLAPI void DrawCubeWires(Vector3 position, float width, float height, float length, Color color)
```

**Fuji usage**

```fuji
let result = DrawCubeWires(position, width, height, length, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `width` | `float` |
| `height` | `float` |
| `length` | `float` |
| `color` | `Color` |

---

### DrawCubeWiresV

```c
RLAPI void DrawCubeWiresV(Vector3 position, Vector3 size, Color color)
```

**Fuji usage**

```fuji
let result = DrawCubeWiresV(position, size, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `size` | `Vector3` |
| `color` | `Color` |

---

### DrawSphere

```c
RLAPI void DrawSphere(Vector3 centerPos, float radius, Color color)
```

**Fuji usage**

```fuji
let result = DrawSphere(centerPos, radius, color);
```

| Parameter | Type |
|-----------|------|
| `centerPos` | `Vector3` |
| `radius` | `float` |
| `color` | `Color` |

---

### DrawSphereEx

```c
RLAPI void DrawSphereEx(Vector3 centerPos, float radius, int rings, int slices, Color color)
```

**Fuji usage**

```fuji
let result = DrawSphereEx(centerPos, radius, rings, slices, color);
```

| Parameter | Type |
|-----------|------|
| `centerPos` | `Vector3` |
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawSphereWires

```c
RLAPI void DrawSphereWires(Vector3 centerPos, float radius, int rings, int slices, Color color)
```

**Fuji usage**

```fuji
let result = DrawSphereWires(centerPos, radius, rings, slices, color);
```

| Parameter | Type |
|-----------|------|
| `centerPos` | `Vector3` |
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCylinder

```c
RLAPI void DrawCylinder(Vector3 position, float radiusTop, float radiusBottom, float height, int slices, Color color)
```

**Fuji usage**

```fuji
let result = DrawCylinder(position, radiusTop, radiusBottom, height, slices, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `radiusTop` | `float` |
| `radiusBottom` | `float` |
| `height` | `float` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCylinderEx

```c
RLAPI void DrawCylinderEx(Vector3 startPos, Vector3 endPos, float startRadius, float endRadius, int sides, Color color)
```

**Fuji usage**

```fuji
let result = DrawCylinderEx(startPos, endPos, startRadius, endRadius, sides, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `startRadius` | `float` |
| `endRadius` | `float` |
| `sides` | `int` |
| `color` | `Color` |

---

### DrawCylinderWires

```c
RLAPI void DrawCylinderWires(Vector3 position, float radiusTop, float radiusBottom, float height, int slices, Color color)
```

**Fuji usage**

```fuji
let result = DrawCylinderWires(position, radiusTop, radiusBottom, height, slices, color);
```

| Parameter | Type |
|-----------|------|
| `position` | `Vector3` |
| `radiusTop` | `float` |
| `radiusBottom` | `float` |
| `height` | `float` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCylinderWiresEx

```c
RLAPI void DrawCylinderWiresEx(Vector3 startPos, Vector3 endPos, float startRadius, float endRadius, int slices, Color color)
```

**Fuji usage**

```fuji
let result = DrawCylinderWiresEx(startPos, endPos, startRadius, endRadius, slices, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `startRadius` | `float` |
| `endRadius` | `float` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCapsule

```c
RLAPI void DrawCapsule(Vector3 startPos, Vector3 endPos, float radius, int rings, int slices, Color color)
```

**Fuji usage**

```fuji
let result = DrawCapsule(startPos, endPos, radius, rings, slices, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawCapsuleWires

```c
RLAPI void DrawCapsuleWires(Vector3 startPos, Vector3 endPos, float radius, int rings, int slices, Color color)
```

**Fuji usage**

```fuji
let result = DrawCapsuleWires(startPos, endPos, radius, rings, slices, color);
```

| Parameter | Type |
|-----------|------|
| `startPos` | `Vector3` |
| `endPos` | `Vector3` |
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |
| `color` | `Color` |

---

### DrawPlane

```c
RLAPI void DrawPlane(Vector3 centerPos, Vector2 size, Color color)
```

**Fuji usage**

```fuji
let result = DrawPlane(centerPos, size, color);
```

| Parameter | Type |
|-----------|------|
| `centerPos` | `Vector3` |
| `size` | `Vector2` |
| `color` | `Color` |

---

### DrawRay

```c
RLAPI void DrawRay(Ray ray, Color color)
```

**Fuji usage**

```fuji
let result = DrawRay(ray, color);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `color` | `Color` |

---

### DrawGrid

```c
RLAPI void DrawGrid(int slices, float spacing)
```

**Fuji usage**

```fuji
let result = DrawGrid(slices, spacing);
```

| Parameter | Type |
|-----------|------|
| `slices` | `int` |
| `spacing` | `float` |

---

### LoadModel

```c
RLAPI Model LoadModel(const char *fileName)
```

**Fuji usage**

```fuji
let result = LoadModel(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### LoadModelFromMesh

```c
RLAPI Model LoadModelFromMesh(Mesh mesh)
```

**Fuji usage**

```fuji
let result = LoadModelFromMesh(mesh);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |

---

### IsModelValid

```c
RLAPI bool IsModelValid(Model model)
```

**Fuji usage**

```fuji
let result = IsModelValid(model);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |

---

### UnloadModel

```c
RLAPI void UnloadModel(Model model)
```

**Fuji usage**

```fuji
let result = UnloadModel(model);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |

---

### GetModelBoundingBox

```c
RLAPI BoundingBox GetModelBoundingBox(Model model)
```

**Fuji usage**

```fuji
let result = GetModelBoundingBox(model);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |

---

### DrawModel

```c
RLAPI void DrawModel(Model model, Vector3 position, float scale, Color tint)
```

**Fuji usage**

```fuji
let result = DrawModel(model, position, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `position` | `Vector3` |
| `scale` | `float` |
| `tint` | `Color` |

---

### DrawModelEx

```c
RLAPI void DrawModelEx(Model model, Vector3 position, Vector3 rotationAxis, float rotationAngle, Vector3 scale, Color tint)
```

**Fuji usage**

```fuji
let result = DrawModelEx(model, position, rotationAxis, rotationAngle, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `position` | `Vector3` |
| `rotationAxis` | `Vector3` |
| `rotationAngle` | `float` |
| `scale` | `Vector3` |
| `tint` | `Color` |

---

### DrawModelWires

```c
RLAPI void DrawModelWires(Model model, Vector3 position, float scale, Color tint)
```

**Fuji usage**

```fuji
let result = DrawModelWires(model, position, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `position` | `Vector3` |
| `scale` | `float` |
| `tint` | `Color` |

---

### DrawModelWiresEx

```c
RLAPI void DrawModelWiresEx(Model model, Vector3 position, Vector3 rotationAxis, float rotationAngle, Vector3 scale, Color tint)
```

**Fuji usage**

```fuji
let result = DrawModelWiresEx(model, position, rotationAxis, rotationAngle, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `position` | `Vector3` |
| `rotationAxis` | `Vector3` |
| `rotationAngle` | `float` |
| `scale` | `Vector3` |
| `tint` | `Color` |

---

### DrawBoundingBox

```c
RLAPI void DrawBoundingBox(BoundingBox box, Color color)
```

**Fuji usage**

```fuji
let result = DrawBoundingBox(box, color);
```

| Parameter | Type |
|-----------|------|
| `box` | `BoundingBox` |
| `color` | `Color` |

---

### DrawBillboard

```c
RLAPI void DrawBillboard(Camera camera, Texture2D texture, Vector3 position, float scale, Color tint)
```

**Fuji usage**

```fuji
let result = DrawBillboard(camera, texture, position, scale, tint);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera` |
| `texture` | `Texture2D` |
| `position` | `Vector3` |
| `scale` | `float` |
| `tint` | `Color` |

---

### DrawBillboardRec

```c
RLAPI void DrawBillboardRec(Camera camera, Texture2D texture, Rectangle source, Vector3 position, Vector2 size, Color tint)
```

**Fuji usage**

```fuji
let result = DrawBillboardRec(camera, texture, source, position, size, tint);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera` |
| `texture` | `Texture2D` |
| `source` | `Rectangle` |
| `position` | `Vector3` |
| `size` | `Vector2` |
| `tint` | `Color` |

---

### DrawBillboardPro

```c
RLAPI void DrawBillboardPro(Camera camera, Texture2D texture, Rectangle source, Vector3 position, Vector3 up, Vector2 size, Vector2 origin, float rotation, Color tint)
```

**Fuji usage**

```fuji
let result = DrawBillboardPro(camera, texture, source, position, up, size, origin, rotation, tint);
```

| Parameter | Type |
|-----------|------|
| `camera` | `Camera` |
| `texture` | `Texture2D` |
| `source` | `Rectangle` |
| `position` | `Vector3` |
| `up` | `Vector3` |
| `size` | `Vector2` |
| `origin` | `Vector2` |
| `rotation` | `float` |
| `tint` | `Color` |

---

### UploadMesh

```c
RLAPI void UploadMesh(Mesh *mesh, bool dynamic)
```

**Fuji usage**

```fuji
let result = UploadMesh(*mesh, dynamic);
```

| Parameter | Type |
|-----------|------|
| `*mesh` | `Mesh` |
| `dynamic` | `bool` |

---

### UpdateMeshBuffer

```c
RLAPI void UpdateMeshBuffer(Mesh mesh, int index, const void *data, int dataSize, int offset)
```

**Fuji usage**

```fuji
let result = UpdateMeshBuffer(mesh, index, *data, dataSize, offset);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `index` | `int` |
| `*data` | `const void` |
| `dataSize` | `int` |
| `offset` | `int` |

---

### UnloadMesh

```c
RLAPI void UnloadMesh(Mesh mesh)
```

**Fuji usage**

```fuji
let result = UnloadMesh(mesh);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |

---

### DrawMesh

```c
RLAPI void DrawMesh(Mesh mesh, Material material, Matrix transform)
```

**Fuji usage**

```fuji
let result = DrawMesh(mesh, material, transform);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `material` | `Material` |
| `transform` | `Matrix` |

---

### DrawMeshInstanced

```c
RLAPI void DrawMeshInstanced(Mesh mesh, Material material, const Matrix *transforms, int instances)
```

**Fuji usage**

```fuji
let result = DrawMeshInstanced(mesh, material, *transforms, instances);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `material` | `Material` |
| `*transforms` | `const Matrix` |
| `instances` | `int` |

---

### GetMeshBoundingBox

```c
RLAPI BoundingBox GetMeshBoundingBox(Mesh mesh)
```

**Fuji usage**

```fuji
let result = GetMeshBoundingBox(mesh);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |

---

### GenMeshTangents

```c
RLAPI void GenMeshTangents(Mesh *mesh)
```

**Fuji usage**

```fuji
let result = GenMeshTangents(*mesh);
```

| Parameter | Type |
|-----------|------|
| `*mesh` | `Mesh` |

---

### ExportMesh

```c
RLAPI bool ExportMesh(Mesh mesh, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportMesh(mesh, *fileName);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `*fileName` | `const char` |

---

### ExportMeshAsCode

```c
RLAPI bool ExportMeshAsCode(Mesh mesh, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportMeshAsCode(mesh, *fileName);
```

| Parameter | Type |
|-----------|------|
| `mesh` | `Mesh` |
| `*fileName` | `const char` |

---

### GenMeshPoly

```c
RLAPI Mesh GenMeshPoly(int sides, float radius)
```

**Fuji usage**

```fuji
let result = GenMeshPoly(sides, radius);
```

| Parameter | Type |
|-----------|------|
| `sides` | `int` |
| `radius` | `float` |

---

### GenMeshPlane

```c
RLAPI Mesh GenMeshPlane(float width, float length, int resX, int resZ)
```

**Fuji usage**

```fuji
let result = GenMeshPlane(width, length, resX, resZ);
```

| Parameter | Type |
|-----------|------|
| `width` | `float` |
| `length` | `float` |
| `resX` | `int` |
| `resZ` | `int` |

---

### GenMeshCube

```c
RLAPI Mesh GenMeshCube(float width, float height, float length)
```

**Fuji usage**

```fuji
let result = GenMeshCube(width, height, length);
```

| Parameter | Type |
|-----------|------|
| `width` | `float` |
| `height` | `float` |
| `length` | `float` |

---

### GenMeshSphere

```c
RLAPI Mesh GenMeshSphere(float radius, int rings, int slices)
```

**Fuji usage**

```fuji
let result = GenMeshSphere(radius, rings, slices);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |

---

### GenMeshHemiSphere

```c
RLAPI Mesh GenMeshHemiSphere(float radius, int rings, int slices)
```

**Fuji usage**

```fuji
let result = GenMeshHemiSphere(radius, rings, slices);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `rings` | `int` |
| `slices` | `int` |

---

### GenMeshCylinder

```c
RLAPI Mesh GenMeshCylinder(float radius, float height, int slices)
```

**Fuji usage**

```fuji
let result = GenMeshCylinder(radius, height, slices);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `height` | `float` |
| `slices` | `int` |

---

### GenMeshCone

```c
RLAPI Mesh GenMeshCone(float radius, float height, int slices)
```

**Fuji usage**

```fuji
let result = GenMeshCone(radius, height, slices);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `height` | `float` |
| `slices` | `int` |

---

### GenMeshTorus

```c
RLAPI Mesh GenMeshTorus(float radius, float size, int radSeg, int sides)
```

**Fuji usage**

```fuji
let result = GenMeshTorus(radius, size, radSeg, sides);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `size` | `float` |
| `radSeg` | `int` |
| `sides` | `int` |

---

### GenMeshKnot

```c
RLAPI Mesh GenMeshKnot(float radius, float size, int radSeg, int sides)
```

**Fuji usage**

```fuji
let result = GenMeshKnot(radius, size, radSeg, sides);
```

| Parameter | Type |
|-----------|------|
| `radius` | `float` |
| `size` | `float` |
| `radSeg` | `int` |
| `sides` | `int` |

---

### GenMeshHeightmap

```c
RLAPI Mesh GenMeshHeightmap(Image heightmap, Vector3 size)
```

**Fuji usage**

```fuji
let result = GenMeshHeightmap(heightmap, size);
```

| Parameter | Type |
|-----------|------|
| `heightmap` | `Image` |
| `size` | `Vector3` |

---

### GenMeshCubicmap

```c
RLAPI Mesh GenMeshCubicmap(Image cubicmap, Vector3 cubeSize)
```

**Fuji usage**

```fuji
let result = GenMeshCubicmap(cubicmap, cubeSize);
```

| Parameter | Type |
|-----------|------|
| `cubicmap` | `Image` |
| `cubeSize` | `Vector3` |

---

### LoadMaterialDefault

```c
RLAPI Material LoadMaterialDefault()
```

**Fuji usage**

```fuji
let result = LoadMaterialDefault();
```

---

### IsMaterialValid

```c
RLAPI bool IsMaterialValid(Material material)
```

**Fuji usage**

```fuji
let result = IsMaterialValid(material);
```

| Parameter | Type |
|-----------|------|
| `material` | `Material` |

---

### UnloadMaterial

```c
RLAPI void UnloadMaterial(Material material)
```

**Fuji usage**

```fuji
let result = UnloadMaterial(material);
```

| Parameter | Type |
|-----------|------|
| `material` | `Material` |

---

### SetMaterialTexture

```c
RLAPI void SetMaterialTexture(Material *material, int mapType, Texture2D texture)
```

**Fuji usage**

```fuji
let result = SetMaterialTexture(*material, mapType, texture);
```

| Parameter | Type |
|-----------|------|
| `*material` | `Material` |
| `mapType` | `int` |
| `texture` | `Texture2D` |

---

### SetModelMeshMaterial

```c
RLAPI void SetModelMeshMaterial(Model *model, int meshId, int materialId)
```

**Fuji usage**

```fuji
let result = SetModelMeshMaterial(*model, meshId, materialId);
```

| Parameter | Type |
|-----------|------|
| `*model` | `Model` |
| `meshId` | `int` |
| `materialId` | `int` |

---

### UpdateModelAnimation

```c
RLAPI void UpdateModelAnimation(Model model, ModelAnimation anim, float frame)
```

**Fuji usage**

```fuji
let result = UpdateModelAnimation(model, anim, frame);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `anim` | `ModelAnimation` |
| `frame` | `float` |

---

### UpdateModelAnimationEx

```c
RLAPI void UpdateModelAnimationEx(Model model, ModelAnimation animA, float frameA, ModelAnimation animB, float frameB, float blend)
```

**Fuji usage**

```fuji
let result = UpdateModelAnimationEx(model, animA, frameA, animB, frameB, blend);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `animA` | `ModelAnimation` |
| `frameA` | `float` |
| `animB` | `ModelAnimation` |
| `frameB` | `float` |
| `blend` | `float` |

---

### UnloadModelAnimations

```c
RLAPI void UnloadModelAnimations(ModelAnimation *animations, int animCount)
```

**Fuji usage**

```fuji
let result = UnloadModelAnimations(*animations, animCount);
```

| Parameter | Type |
|-----------|------|
| `*animations` | `ModelAnimation` |
| `animCount` | `int` |

---

### IsModelAnimationValid

```c
RLAPI bool IsModelAnimationValid(Model model, ModelAnimation anim)
```

**Fuji usage**

```fuji
let result = IsModelAnimationValid(model, anim);
```

| Parameter | Type |
|-----------|------|
| `model` | `Model` |
| `anim` | `ModelAnimation` |

---

### CheckCollisionSpheres

```c
RLAPI bool CheckCollisionSpheres(Vector3 center1, float radius1, Vector3 center2, float radius2)
```

**Fuji usage**

```fuji
let result = CheckCollisionSpheres(center1, radius1, center2, radius2);
```

| Parameter | Type |
|-----------|------|
| `center1` | `Vector3` |
| `radius1` | `float` |
| `center2` | `Vector3` |
| `radius2` | `float` |

---

### CheckCollisionBoxes

```c
RLAPI bool CheckCollisionBoxes(BoundingBox box1, BoundingBox box2)
```

**Fuji usage**

```fuji
let result = CheckCollisionBoxes(box1, box2);
```

| Parameter | Type |
|-----------|------|
| `box1` | `BoundingBox` |
| `box2` | `BoundingBox` |

---

### CheckCollisionBoxSphere

```c
RLAPI bool CheckCollisionBoxSphere(BoundingBox box, Vector3 center, float radius)
```

**Fuji usage**

```fuji
let result = CheckCollisionBoxSphere(box, center, radius);
```

| Parameter | Type |
|-----------|------|
| `box` | `BoundingBox` |
| `center` | `Vector3` |
| `radius` | `float` |

---

### GetRayCollisionSphere

```c
RLAPI RayCollision GetRayCollisionSphere(Ray ray, Vector3 center, float radius)
```

**Fuji usage**

```fuji
let result = GetRayCollisionSphere(ray, center, radius);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `center` | `Vector3` |
| `radius` | `float` |

---

### GetRayCollisionBox

```c
RLAPI RayCollision GetRayCollisionBox(Ray ray, BoundingBox box)
```

**Fuji usage**

```fuji
let result = GetRayCollisionBox(ray, box);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `box` | `BoundingBox` |

---

### GetRayCollisionMesh

```c
RLAPI RayCollision GetRayCollisionMesh(Ray ray, Mesh mesh, Matrix transform)
```

**Fuji usage**

```fuji
let result = GetRayCollisionMesh(ray, mesh, transform);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `mesh` | `Mesh` |
| `transform` | `Matrix` |

---

### GetRayCollisionTriangle

```c
RLAPI RayCollision GetRayCollisionTriangle(Ray ray, Vector3 p1, Vector3 p2, Vector3 p3)
```

**Fuji usage**

```fuji
let result = GetRayCollisionTriangle(ray, p1, p2, p3);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `p1` | `Vector3` |
| `p2` | `Vector3` |
| `p3` | `Vector3` |

---

### GetRayCollisionQuad

```c
RLAPI RayCollision GetRayCollisionQuad(Ray ray, Vector3 p1, Vector3 p2, Vector3 p3, Vector3 p4)
```

**Fuji usage**

```fuji
let result = GetRayCollisionQuad(ray, p1, p2, p3, p4);
```

| Parameter | Type |
|-----------|------|
| `ray` | `Ray` |
| `p1` | `Vector3` |
| `p2` | `Vector3` |
| `p3` | `Vector3` |
| `p4` | `Vector3` |

---

### void

```c
typedef void()
```

**Fuji usage**

```fuji
let result = void();
```

---

### InitAudioDevice

```c
RLAPI void InitAudioDevice()
```

**Fuji usage**

```fuji
let result = InitAudioDevice();
```

---

### CloseAudioDevice

```c
RLAPI void CloseAudioDevice()
```

**Fuji usage**

```fuji
let result = CloseAudioDevice();
```

---

### IsAudioDeviceReady

```c
RLAPI bool IsAudioDeviceReady()
```

**Fuji usage**

```fuji
let result = IsAudioDeviceReady();
```

---

### SetMasterVolume

```c
RLAPI void SetMasterVolume(float volume)
```

**Fuji usage**

```fuji
let result = SetMasterVolume(volume);
```

| Parameter | Type |
|-----------|------|
| `volume` | `float` |

---

### GetMasterVolume

```c
RLAPI float GetMasterVolume()
```

**Fuji usage**

```fuji
let result = GetMasterVolume();
```

---

### LoadWave

```c
RLAPI Wave LoadWave(const char *fileName)
```

**Fuji usage**

```fuji
let result = LoadWave(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### LoadWaveFromMemory

```c
RLAPI Wave LoadWaveFromMemory(const char *fileType, const unsigned char *fileData, int dataSize)
```

**Fuji usage**

```fuji
let result = LoadWaveFromMemory(*fileType, *fileData, dataSize);
```

| Parameter | Type |
|-----------|------|
| `*fileType` | `const char` |
| `*fileData` | `const unsigned char` |
| `dataSize` | `int` |

---

### IsWaveValid

```c
RLAPI bool IsWaveValid(Wave wave)
```

**Fuji usage**

```fuji
let result = IsWaveValid(wave);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |

---

### LoadSound

```c
RLAPI Sound LoadSound(const char *fileName)
```

**Fuji usage**

```fuji
let result = LoadSound(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### LoadSoundFromWave

```c
RLAPI Sound LoadSoundFromWave(Wave wave)
```

**Fuji usage**

```fuji
let result = LoadSoundFromWave(wave);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |

---

### LoadSoundAlias

```c
RLAPI Sound LoadSoundAlias(Sound source)
```

**Fuji usage**

```fuji
let result = LoadSoundAlias(source);
```

| Parameter | Type |
|-----------|------|
| `source` | `Sound` |

---

### IsSoundValid

```c
RLAPI bool IsSoundValid(Sound sound)
```

**Fuji usage**

```fuji
let result = IsSoundValid(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### UpdateSound

```c
RLAPI void UpdateSound(Sound sound, const void *data, int frameCount)
```

**Fuji usage**

```fuji
let result = UpdateSound(sound, *data, frameCount);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |
| `*data` | `const void` |
| `frameCount` | `int` |

---

### UnloadWave

```c
RLAPI void UnloadWave(Wave wave)
```

**Fuji usage**

```fuji
let result = UnloadWave(wave);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |

---

### UnloadSound

```c
RLAPI void UnloadSound(Sound sound)
```

**Fuji usage**

```fuji
let result = UnloadSound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### UnloadSoundAlias

```c
RLAPI void UnloadSoundAlias(Sound alias)
```

**Fuji usage**

```fuji
let result = UnloadSoundAlias(alias);
```

| Parameter | Type |
|-----------|------|
| `alias` | `Sound` |

---

### ExportWave

```c
RLAPI bool ExportWave(Wave wave, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportWave(wave, *fileName);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |
| `*fileName` | `const char` |

---

### ExportWaveAsCode

```c
RLAPI bool ExportWaveAsCode(Wave wave, const char *fileName)
```

**Fuji usage**

```fuji
let result = ExportWaveAsCode(wave, *fileName);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |
| `*fileName` | `const char` |

---

### PlaySound

```c
RLAPI void PlaySound(Sound sound)
```

**Fuji usage**

```fuji
let result = PlaySound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### StopSound

```c
RLAPI void StopSound(Sound sound)
```

**Fuji usage**

```fuji
let result = StopSound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### PauseSound

```c
RLAPI void PauseSound(Sound sound)
```

**Fuji usage**

```fuji
let result = PauseSound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### ResumeSound

```c
RLAPI void ResumeSound(Sound sound)
```

**Fuji usage**

```fuji
let result = ResumeSound(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### IsSoundPlaying

```c
RLAPI bool IsSoundPlaying(Sound sound)
```

**Fuji usage**

```fuji
let result = IsSoundPlaying(sound);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |

---

### SetSoundVolume

```c
RLAPI void SetSoundVolume(Sound sound, float volume)
```

**Fuji usage**

```fuji
let result = SetSoundVolume(sound, volume);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |
| `volume` | `float` |

---

### SetSoundPitch

```c
RLAPI void SetSoundPitch(Sound sound, float pitch)
```

**Fuji usage**

```fuji
let result = SetSoundPitch(sound, pitch);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |
| `pitch` | `float` |

---

### SetSoundPan

```c
RLAPI void SetSoundPan(Sound sound, float pan)
```

**Fuji usage**

```fuji
let result = SetSoundPan(sound, pan);
```

| Parameter | Type |
|-----------|------|
| `sound` | `Sound` |
| `pan` | `float` |

---

### WaveCopy

```c
RLAPI Wave WaveCopy(Wave wave)
```

**Fuji usage**

```fuji
let result = WaveCopy(wave);
```

| Parameter | Type |
|-----------|------|
| `wave` | `Wave` |

---

### WaveCrop

```c
RLAPI void WaveCrop(Wave *wave, int initFrame, int finalFrame)
```

**Fuji usage**

```fuji
let result = WaveCrop(*wave, initFrame, finalFrame);
```

| Parameter | Type |
|-----------|------|
| `*wave` | `Wave` |
| `initFrame` | `int` |
| `finalFrame` | `int` |

---

### WaveFormat

```c
RLAPI void WaveFormat(Wave *wave, int sampleRate, int sampleSize, int channels)
```

**Fuji usage**

```fuji
let result = WaveFormat(*wave, sampleRate, sampleSize, channels);
```

| Parameter | Type |
|-----------|------|
| `*wave` | `Wave` |
| `sampleRate` | `int` |
| `sampleSize` | `int` |
| `channels` | `int` |

---

### UnloadWaveSamples

```c
RLAPI void UnloadWaveSamples(float *samples)
```

**Fuji usage**

```fuji
let result = UnloadWaveSamples(*samples);
```

| Parameter | Type |
|-----------|------|
| `*samples` | `float` |

---

### LoadMusicStream

```c
RLAPI Music LoadMusicStream(const char *fileName)
```

**Fuji usage**

```fuji
let result = LoadMusicStream(*fileName);
```

| Parameter | Type |
|-----------|------|
| `*fileName` | `const char` |

---

### LoadMusicStreamFromMemory

```c
RLAPI Music LoadMusicStreamFromMemory(const char *fileType, const unsigned char *data, int dataSize)
```

**Fuji usage**

```fuji
let result = LoadMusicStreamFromMemory(*fileType, *data, dataSize);
```

| Parameter | Type |
|-----------|------|
| `*fileType` | `const char` |
| `*data` | `const unsigned char` |
| `dataSize` | `int` |

---

### IsMusicValid

```c
RLAPI bool IsMusicValid(Music music)
```

**Fuji usage**

```fuji
let result = IsMusicValid(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### UnloadMusicStream

```c
RLAPI void UnloadMusicStream(Music music)
```

**Fuji usage**

```fuji
let result = UnloadMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### PlayMusicStream

```c
RLAPI void PlayMusicStream(Music music)
```

**Fuji usage**

```fuji
let result = PlayMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### IsMusicStreamPlaying

```c
RLAPI bool IsMusicStreamPlaying(Music music)
```

**Fuji usage**

```fuji
let result = IsMusicStreamPlaying(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### UpdateMusicStream

```c
RLAPI void UpdateMusicStream(Music music)
```

**Fuji usage**

```fuji
let result = UpdateMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### StopMusicStream

```c
RLAPI void StopMusicStream(Music music)
```

**Fuji usage**

```fuji
let result = StopMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### PauseMusicStream

```c
RLAPI void PauseMusicStream(Music music)
```

**Fuji usage**

```fuji
let result = PauseMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### ResumeMusicStream

```c
RLAPI void ResumeMusicStream(Music music)
```

**Fuji usage**

```fuji
let result = ResumeMusicStream(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### SeekMusicStream

```c
RLAPI void SeekMusicStream(Music music, float position)
```

**Fuji usage**

```fuji
let result = SeekMusicStream(music, position);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |
| `position` | `float` |

---

### SetMusicVolume

```c
RLAPI void SetMusicVolume(Music music, float volume)
```

**Fuji usage**

```fuji
let result = SetMusicVolume(music, volume);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |
| `volume` | `float` |

---

### SetMusicPitch

```c
RLAPI void SetMusicPitch(Music music, float pitch)
```

**Fuji usage**

```fuji
let result = SetMusicPitch(music, pitch);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |
| `pitch` | `float` |

---

### SetMusicPan

```c
RLAPI void SetMusicPan(Music music, float pan)
```

**Fuji usage**

```fuji
let result = SetMusicPan(music, pan);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |
| `pan` | `float` |

---

### GetMusicTimeLength

```c
RLAPI float GetMusicTimeLength(Music music)
```

**Fuji usage**

```fuji
let result = GetMusicTimeLength(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### GetMusicTimePlayed

```c
RLAPI float GetMusicTimePlayed(Music music)
```

**Fuji usage**

```fuji
let result = GetMusicTimePlayed(music);
```

| Parameter | Type |
|-----------|------|
| `music` | `Music` |

---

### LoadAudioStream

```c
RLAPI AudioStream LoadAudioStream(unsigned int sampleRate, unsigned int sampleSize, unsigned int channels)
```

**Fuji usage**

```fuji
let result = LoadAudioStream(sampleRate, sampleSize, channels);
```

| Parameter | Type |
|-----------|------|
| `sampleRate` | `unsigned int` |
| `sampleSize` | `unsigned int` |
| `channels` | `unsigned int` |

---

### IsAudioStreamValid

```c
RLAPI bool IsAudioStreamValid(AudioStream stream)
```

**Fuji usage**

```fuji
let result = IsAudioStreamValid(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### UnloadAudioStream

```c
RLAPI void UnloadAudioStream(AudioStream stream)
```

**Fuji usage**

```fuji
let result = UnloadAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### UpdateAudioStream

```c
RLAPI void UpdateAudioStream(AudioStream stream, const void *data, int frameCount)
```

**Fuji usage**

```fuji
let result = UpdateAudioStream(stream, *data, frameCount);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `*data` | `const void` |
| `frameCount` | `int` |

---

### IsAudioStreamProcessed

```c
RLAPI bool IsAudioStreamProcessed(AudioStream stream)
```

**Fuji usage**

```fuji
let result = IsAudioStreamProcessed(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### PlayAudioStream

```c
RLAPI void PlayAudioStream(AudioStream stream)
```

**Fuji usage**

```fuji
let result = PlayAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### PauseAudioStream

```c
RLAPI void PauseAudioStream(AudioStream stream)
```

**Fuji usage**

```fuji
let result = PauseAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### ResumeAudioStream

```c
RLAPI void ResumeAudioStream(AudioStream stream)
```

**Fuji usage**

```fuji
let result = ResumeAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### IsAudioStreamPlaying

```c
RLAPI bool IsAudioStreamPlaying(AudioStream stream)
```

**Fuji usage**

```fuji
let result = IsAudioStreamPlaying(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### StopAudioStream

```c
RLAPI void StopAudioStream(AudioStream stream)
```

**Fuji usage**

```fuji
let result = StopAudioStream(stream);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |

---

### SetAudioStreamVolume

```c
RLAPI void SetAudioStreamVolume(AudioStream stream, float volume)
```

**Fuji usage**

```fuji
let result = SetAudioStreamVolume(stream, volume);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `volume` | `float` |

---

### SetAudioStreamPitch

```c
RLAPI void SetAudioStreamPitch(AudioStream stream, float pitch)
```

**Fuji usage**

```fuji
let result = SetAudioStreamPitch(stream, pitch);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `pitch` | `float` |

---

### SetAudioStreamPan

```c
RLAPI void SetAudioStreamPan(AudioStream stream, float pan)
```

**Fuji usage**

```fuji
let result = SetAudioStreamPan(stream, pan);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `pan` | `float` |

---

### SetAudioStreamBufferSizeDefault

```c
RLAPI void SetAudioStreamBufferSizeDefault(int size)
```

**Fuji usage**

```fuji
let result = SetAudioStreamBufferSizeDefault(size);
```

| Parameter | Type |
|-----------|------|
| `size` | `int` |

---

### SetAudioStreamCallback

```c
RLAPI void SetAudioStreamCallback(AudioStream stream, AudioCallback callback)
```

**Fuji usage**

```fuji
let result = SetAudioStreamCallback(stream, callback);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `callback` | `AudioCallback` |

---

### AttachAudioStreamProcessor

```c
RLAPI void AttachAudioStreamProcessor(AudioStream stream, AudioCallback processor)
```

**Fuji usage**

```fuji
let result = AttachAudioStreamProcessor(stream, processor);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `processor` | `AudioCallback` |

---

### DetachAudioStreamProcessor

```c
RLAPI void DetachAudioStreamProcessor(AudioStream stream, AudioCallback processor)
```

**Fuji usage**

```fuji
let result = DetachAudioStreamProcessor(stream, processor);
```

| Parameter | Type |
|-----------|------|
| `stream` | `AudioStream` |
| `processor` | `AudioCallback` |

---

### AttachAudioMixedProcessor

```c
RLAPI void AttachAudioMixedProcessor(AudioCallback processor)
```

**Fuji usage**

```fuji
let result = AttachAudioMixedProcessor(processor);
```

| Parameter | Type |
|-----------|------|
| `processor` | `AudioCallback` |

---

### DetachAudioMixedProcessor

```c
RLAPI void DetachAudioMixedProcessor(AudioCallback processor)
```

**Fuji usage**

```fuji
let result = DetachAudioMixedProcessor(processor);
```

| Parameter | Type |
|-----------|------|
| `processor` | `AudioCallback` |

---

## Structs

### Vector2

| Field | Type |
|-------|------|
| `x` | `float` |
| `y` | `// Vector x component float` |
| `component` | `// Vector y` |

---

### Vector3

| Field | Type |
|-------|------|
| `x` | `float` |
| `y` | `// Vector x component float` |
| `z` | `// Vector y component float` |
| `component` | `// Vector z` |

---

### Vector4

| Field | Type |
|-------|------|
| `x` | `float` |
| `y` | `// Vector x component float` |
| `z` | `// Vector y component float` |
| `w` | `// Vector z component float` |
| `component` | `// Vector w` |

---

### Matrix

| Field | Type |
|-------|------|
| `m12` | `float m0, m4, m8,` |
| `m13` | `// Matrix first row (4 components) float m1, m5, m9,` |
| `m14` | `// Matrix second row (4 components) float m2, m6, m10,` |
| `m15` | `// Matrix third row (4 components) float m3, m7, m11,` |
| `components)` | `// Matrix fourth row (4` |

---

### Color

| Field | Type |
|-------|------|
| `r` | `unsigned char` |
| `g` | `// Color red value unsigned char` |
| `b` | `// Color green value unsigned char` |
| `a` | `// Color blue value unsigned char` |
| `value` | `// Color alpha` |

---

### Rectangle

| Field | Type |
|-------|------|
| `x` | `float` |
| `y` | `// Rectangle top-left corner position x float` |
| `width` | `// Rectangle top-left corner position y float` |
| `height` | `// Rectangle width float` |
| `height` | `// Rectangle` |

---

### Image

| Field | Type |
|-------|------|
| `*data` | `void` |
| `width` | `// Image raw data int` |
| `height` | `// Image base width int` |
| `mipmaps` | `// Image base height int` |
| `format` | `// Mipmap levels, 1 by default int` |
| `type)` | `// Data format (PixelFormat` |

---

### Texture

| Field | Type |
|-------|------|
| `id` | `unsigned int` |
| `width` | `// OpenGL texture id int` |
| `height` | `// Texture base width int` |
| `mipmaps` | `// Texture base height int` |
| `format` | `// Mipmap levels, 1 by default int` |
| `type)` | `// Data format (PixelFormat` |

---

### RenderTexture

| Field | Type |
|-------|------|
| `id` | `unsigned int` |
| `texture` | `// OpenGL framebuffer object id Texture` |
| `depth` | `// Color buffer attachment texture Texture` |
| `texture` | `// Depth buffer attachment` |

---

### NPatchInfo

| Field | Type |
|-------|------|
| `source` | `Rectangle` |
| `left` | `// Texture source rectangle int` |
| `top` | `// Left border offset int` |
| `right` | `// Top border offset int` |
| `bottom` | `// Right border offset int` |
| `layout` | `// Bottom border offset int` |
| `3x1` | `// Layout of the n-patch: 3x3, 1x3 or` |

---

### GlyphInfo

| Field | Type |
|-------|------|
| `value` | `int` |
| `offsetX` | `// Character value (Unicode) int` |
| `offsetY` | `// Character offset X when drawing int` |
| `advanceX` | `// Character offset Y when drawing int` |
| `image` | `// Character advance position X Image` |
| `data` | `// Character image` |

---

### Font

| Field | Type |
|-------|------|
| `baseSize` | `int` |
| `glyphCount` | `// Base size (default chars height) int` |
| `glyphPadding` | `// Number of glyph characters int` |
| `texture` | `// Padding around the glyph characters Texture2D` |
| `*recs` | `// Texture atlas containing the glyphs Rectangle` |
| `*glyphs` | `// Rectangles in texture for the glyphs GlyphInfo` |
| `data` | `// Glyphs info` |

---

### Camera3D

| Field | Type |
|-------|------|
| `position` | `Vector3` |
| `target` | `// Camera position Vector3` |
| `up` | `// Camera target it looks-at Vector3` |
| `fovy` | `// Camera up vector (rotation over its axis) float` |
| `projection` | `// Camera field-of-view aperture in Y (degrees) in perspective, used as near plane height in world units in orthographic int` |
| `CAMERA_ORTHOGRAPHIC` | `// Camera projection: CAMERA_PERSPECTIVE or` |

---

### Camera2D

| Field | Type |
|-------|------|
| `offset` | `Vector2` |
| `target` | `// Camera offset (screen space offset from window origin) Vector2` |
| `rotation` | `// Camera target (world space target point that is mapped to screen space offset) float` |
| `zoom` | `// Camera rotation in degrees (pivots around target) float` |
| `scale` | `// Camera zoom (scaling around target), must not be set to 0, set to 1.0f for no` |

---

### Mesh

| Field | Type |
|-------|------|
| `vertexCount` | `int` |
| `triangleCount` | `// Number of vertices stored in arrays int` |
| `*vertices` | `// Number of triangles stored (indexed or not) // Vertex attributes data float` |
| `*texcoords` | `// Vertex position (XYZ - 3 components per vertex) (shader-location = 0) float` |
| `*texcoords2` | `// Vertex texture coordinates (UV - 2 components per vertex) (shader-location = 1) float` |
| `*normals` | `// Vertex texture second coordinates (UV - 2 components per vertex) (shader-location = 5) float` |
| `*tangents` | `// Vertex normals (XYZ - 3 components per vertex) (shader-location = 2) float` |
| `*colors` | `// Vertex tangents (XYZW - 4 components per vertex) (shader-location = 4) unsigned char` |
| `*indices` | `// Vertex colors (RGBA - 4 components per vertex) (shader-location = 3) unsigned short` |
| `boneCount` | `// Vertex indices (in case vertex data comes indexed) // Skin data for animation int` |
| `*boneIndices` | `// Number of bones (MAX: 256 bones) unsigned char` |
| `*boneWeights` | `// Vertex bone indices, up to 4 bones influence by vertex (skinning) (shader-location = 6) float` |
| `*animVertices` | `// Vertex bone weight, up to 4 bones influence by vertex (skinning) (shader-location = 7) // Runtime animation vertex data (CPU skinning) // NOTE: In case of GPU skinning, not used, pointers are NULL float` |
| `*animNormals` | `// Animated vertex positions (after bones transformations) float` |
| `vaoId` | `// Animated normals (after bones transformations) // OpenGL identifiers unsigned int` |
| `*vboId` | `// OpenGL Vertex Array Object id unsigned int` |
| `data)` | `// OpenGL Vertex Buffer Objects id (default vertex` |

---

### Shader

| Field | Type |
|-------|------|
| `id` | `unsigned int` |
| `*locs` | `// Shader program id int` |
| `(RL_MAX_SHADER_LOCATIONS)` | `// Shader locations array` |

---

### MaterialMap

| Field | Type |
|-------|------|
| `texture` | `Texture2D` |
| `color` | `// Material map texture Color` |
| `value` | `// Material map color float` |
| `value` | `// Material map` |

---

### Material

| Field | Type |
|-------|------|
| `shader` | `Shader` |
| `*maps` | `// Material shader MaterialMap` |
| `params[4]` | `// Material maps array (MAX_MATERIAL_MAPS) float` |
| `required)` | `// Material generic parameters (if` |

---

### Transform

| Field | Type |
|-------|------|
| `translation` | `Vector3` |
| `rotation` | `// Translation Quaternion` |
| `scale` | `// Rotation Vector3` |
| `Scale` | `//` |

---

### BoneInfo

| Field | Type |
|-------|------|
| `name[32]` | `char` |
| `parent` | `// Bone name int` |
| `parent` | `// Bone` |

---

### ModelSkeleton

| Field | Type |
|-------|------|
| `boneCount` | `int` |
| `*bones` | `// Number of bones BoneInfo` |
| `bindPose` | `// Bones information (skeleton) ModelAnimPose` |
| `(Transform[])` | `// Bones base transformation` |

---

### Model

| Field | Type |
|-------|------|
| `transform` | `Matrix` |
| `meshCount` | `// Local transform matrix int` |
| `materialCount` | `// Number of meshes int` |
| `*meshes` | `// Number of materials Mesh` |
| `*materials` | `// Meshes array Material` |
| `*meshMaterial` | `// Materials array int` |
| `skeleton` | `// Mesh material number // Animation data ModelSkeleton` |
| `currentPose` | `// Skeleton for animation // Runtime animation data (CPU/GPU skinning) ModelAnimPose` |
| `*boneMatrices` | `// Current animation pose (Transform[]) Matrix` |
| `matrices` | `// Bones animated transformation` |

---

### ModelAnimation

| Field | Type |
|-------|------|
| `name[32]` | `char` |
| `boneCount` | `// Animation name int` |
| `keyframeCount` | `// Number of bones (per pose) int` |
| `*keyframePoses` | `// Number of animation key frames ModelAnimPose` |
| `[keyframe][pose]` | `// Animation sequence keyframe poses` |

---

### Ray

| Field | Type |
|-------|------|
| `position` | `Vector3` |
| `direction` | `// Ray position (origin) Vector3` |
| `(normalized)` | `// Ray direction` |

---

### RayCollision

| Field | Type |
|-------|------|
| `hit` | `bool` |
| `distance` | `// Did the ray hit something? float` |
| `point` | `// Distance to the nearest hit Vector3` |
| `normal` | `// Point of the nearest hit Vector3` |
| `hit` | `// Surface normal of` |

---

### BoundingBox

| Field | Type |
|-------|------|
| `min` | `Vector3` |
| `max` | `// Minimum vertex box-corner Vector3` |
| `box-corner` | `// Maximum vertex` |

---

### Wave

| Field | Type |
|-------|------|
| `frameCount` | `unsigned int` |
| `sampleRate` | `// Total number of frames (considering channels) unsigned int` |
| `sampleSize` | `// Frequency (samples per second) unsigned int` |
| `channels` | `// Bit depth (bits per sample): 8, 16, 32 (24 not supported) unsigned int` |
| `*data` | `// Number of channels (1-mono, 2-stereo, ...) void` |
| `pointer` | `// Buffer data` |

---

### AudioStream

| Field | Type |
|-------|------|
| `*buffer` | `rAudioBuffer` |
| `*processor` | `// Pointer to internal data used by the audio system rAudioProcessor` |
| `sampleRate` | `// Pointer to internal data processor, useful for audio effects unsigned int` |
| `sampleSize` | `// Frequency (samples per second) unsigned int` |
| `channels` | `// Bit depth (bits per sample): 8, 16, 32 (24 not supported) unsigned int` |
| `...)` | `// Number of channels (1-mono, 2-stereo,` |

---

### Sound

| Field | Type |
|-------|------|
| `stream` | `AudioStream` |
| `frameCount` | `// Audio stream unsigned int` |
| `channels)` | `// Total number of frames (considering` |

---

### Music

| Field | Type |
|-------|------|
| `stream` | `AudioStream` |
| `frameCount` | `// Audio stream unsigned int` |
| `looping` | `// Total number of frames (considering channels) bool` |
| `ctxType` | `// Music looping enable int` |
| `*ctxData` | `// Type of music context (audio filetype) void` |
| `type` | `// Audio context data, depends on` |

---

### VrDeviceInfo

| Field | Type |
|-------|------|
| `hResolution` | `int` |
| `vResolution` | `// Horizontal resolution in pixels int` |
| `hScreenSize` | `// Vertical resolution in pixels float` |
| `vScreenSize` | `// Horizontal size in meters float` |
| `eyeToScreenDistance` | `// Vertical size in meters float` |
| `lensSeparationDistance` | `// Distance between eye and display in meters float` |
| `interpupillaryDistance` | `// Lens separation distance in meters float` |
| `lensDistortionValues[4]` | `// IPD (distance between pupils) in meters float` |
| `chromaAbCorrection[4]` | `// Lens distortion constant parameters float` |
| `parameters` | `// Chromatic aberration correction` |

---

### VrStereoConfig

| Field | Type |
|-------|------|
| `projection[2]` | `Matrix` |
| `viewOffset[2]` | `// VR projection matrices (per eye) Matrix` |
| `leftLensCenter[2]` | `// VR view offset matrices (per eye) float` |
| `rightLensCenter[2]` | `// VR left lens center float` |
| `leftScreenCenter[2]` | `// VR right lens center float` |
| `rightScreenCenter[2]` | `// VR left screen center float` |
| `scale[2]` | `// VR right screen center float` |
| `scaleIn[2]` | `// VR distortion scale float` |
| `in` | `// VR distortion scale` |

---

### FilePathList

| Field | Type |
|-------|------|
| `count` | `unsigned int` |
| `**paths` | `// Filepaths entries count char` |
| `entries` | `// Filepaths` |

---

### AutomationEvent

| Field | Type |
|-------|------|
| `frame` | `unsigned int` |
| `type` | `// Event frame unsigned int` |
| `params[4]` | `// Event type (AutomationEventType) int` |
| `required)` | `// Event parameters (if` |

---

### AutomationEventList

| Field | Type |
|-------|------|
| `capacity` | `unsigned int` |
| `count` | `// Events max entries (MAX_AUTOMATION_EVENTS) unsigned int` |
| `*events` | `// Events entries count AutomationEvent` |
| `entries` | `// Events` |

---

## Enums

### bool

| Name | Value |
|------|-------|
| `false` | `0` |
| `true` | `1` |

---

### ConfigFlags

| Name | Value |
|------|-------|
| `FLAG_VSYNC_HINT` | `0` |
| `// Set to try enabling V-Sync on GPU
    FLAG_FULLSCREEN_MODE` | `0` |
| `// Set to run program in fullscreen
    FLAG_WINDOW_RESIZABLE` | `0` |
| `// Set to allow resizable window
    FLAG_WINDOW_UNDECORATED` | `0` |
| `// Set to disable window decoration (frame and buttons)
    FLAG_WINDOW_HIDDEN` | `0` |
| `// Set to hide window
    FLAG_WINDOW_MINIMIZED` | `0` |
| `// Set to minimize window (iconify)
    FLAG_WINDOW_MAXIMIZED` | `0` |
| `// Set to maximize window (expanded to monitor)
    FLAG_WINDOW_UNFOCUSED` | `0` |
| `// Set to window non focused
    FLAG_WINDOW_TOPMOST` | `0` |
| `// Set to window always on top
    FLAG_WINDOW_ALWAYS_RUN` | `0` |
| `// Set to allow windows running while minimized
    FLAG_WINDOW_TRANSPARENT` | `0` |
| `// Set to allow transparent framebuffer
    FLAG_WINDOW_HIGHDPI` | `0` |
| `// Set to support HighDPI
    FLAG_WINDOW_MOUSE_PASSTHROUGH` | `0` |
| `// Set to support mouse passthrough` | `13` |
| `only supported when FLAG_WINDOW_UNDECORATED
    FLAG_BORDERLESS_WINDOWED_MODE` | `0` |
| `// Set to run program in borderless windowed mode
    FLAG_MSAA_4X_HINT` | `0` |
| `// Set to try enabling MSAA 4X
    FLAG_INTERLACED_HINT` | `0` |

---

### TraceLogLevel

| Name | Value |
|------|-------|
| `LOG_ALL` | `0` |
| `// Display all logs
    LOG_TRACE` | `1` |
| `// Trace logging` | `2` |
| `intended for internal use only
    LOG_DEBUG` | `3` |
| `// Debug logging` | `4` |
| `used for internal debugging` | `5` |
| `it should be disabled on release builds
    LOG_INFO` | `6` |
| `// Info logging` | `7` |
| `used for program execution info
    LOG_WARNING` | `8` |
| `// Warning logging` | `9` |
| `used on recoverable failures
    LOG_ERROR` | `10` |
| `// Error logging` | `11` |
| `used on unrecoverable failures
    LOG_FATAL` | `12` |
| `// Fatal logging` | `13` |
| `used to abort program: exit(EXIT_FAILURE)
    LOG_NONE            // Disable logging` | `14` |

---

### KeyboardKey

| Name | Value |
|------|-------|
| `KEY_NULL` | `0` |
| `// Key: NULL` | `1` |
| `used for no key pressed
    // Alphanumeric keys
    KEY_APOSTROPHE` | `39` |
| `// Key: '
    KEY_COMMA` | `44` |
| `// Key:` | `4` |
| `KEY_MINUS` | `45` |
| `// Key: -
    KEY_PERIOD` | `46` |
| `// Key: .
    KEY_SLASH` | `47` |
| `// Key: /
    KEY_ZERO` | `48` |
| `// Key: 0
    KEY_ONE` | `49` |
| `// Key: 1
    KEY_TWO` | `50` |
| `// Key: 2
    KEY_THREE` | `51` |
| `// Key: 3
    KEY_FOUR` | `52` |
| `// Key: 4
    KEY_FIVE` | `53` |
| `// Key: 5
    KEY_SIX` | `54` |
| `// Key: 6
    KEY_SEVEN` | `55` |
| `// Key: 7
    KEY_EIGHT` | `56` |
| `// Key: 8
    KEY_NINE` | `57` |
| `// Key: 9
    KEY_SEMICOLON` | `59` |
| `// Key: ;
    KEY_EQUAL` | `61` |
| `// Key:` | `20` |
| `// Key: A | a
    KEY_B` | `66` |
| `// Key: B | b
    KEY_C` | `67` |
| `// Key: C | c
    KEY_D` | `68` |
| `// Key: D | d
    KEY_E` | `69` |
| `// Key: E | e
    KEY_F` | `70` |
| `// Key: F | f
    KEY_G` | `71` |
| `// Key: G | g
    KEY_H` | `72` |
| `// Key: H | h
    KEY_I` | `73` |
| `// Key: I | i
    KEY_J` | `74` |
| `// Key: J | j
    KEY_K` | `75` |
| `// Key: K | k
    KEY_L` | `76` |
| `// Key: L | l
    KEY_M` | `77` |
| `// Key: M | m
    KEY_N` | `78` |
| `// Key: N | n
    KEY_O` | `79` |
| `// Key: O | o
    KEY_P` | `80` |
| `// Key: P | p
    KEY_Q` | `81` |
| `// Key: Q | q
    KEY_R` | `82` |
| `// Key: R | r
    KEY_S` | `83` |
| `// Key: S | s
    KEY_T` | `84` |
| `// Key: T | t
    KEY_U` | `85` |
| `// Key: U | u
    KEY_V` | `86` |
| `// Key: V | v
    KEY_W` | `87` |
| `// Key: W | w
    KEY_X` | `88` |
| `// Key: X | x
    KEY_Y` | `89` |
| `// Key: Y | y
    KEY_Z` | `90` |
| `// Key: Z | z
    KEY_LEFT_BRACKET` | `91` |
| `// Key: [
    KEY_BACKSLASH` | `92` |
| `// Key: '\'
    KEY_RIGHT_BRACKET` | `93` |
| `// Key: ]
    KEY_GRAVE` | `96` |
| `// Key: `
    // Function keys
    KEY_SPACE` | `32` |
| `// Key: Space
    KEY_ESCAPE` | `256` |
| `// Key: Esc
    KEY_ENTER` | `257` |
| `// Key: Enter
    KEY_TAB` | `258` |
| `// Key: Tab
    KEY_BACKSPACE` | `259` |
| `// Key: Backspace
    KEY_INSERT` | `260` |
| `// Key: Ins
    KEY_DELETE` | `261` |
| `// Key: Del
    KEY_RIGHT` | `262` |
| `// Key: Cursor right
    KEY_LEFT` | `263` |
| `// Key: Cursor left
    KEY_DOWN` | `264` |
| `// Key: Cursor down
    KEY_UP` | `265` |
| `// Key: Cursor up
    KEY_PAGE_UP` | `266` |
| `// Key: Page up
    KEY_PAGE_DOWN` | `267` |
| `// Key: Page down
    KEY_HOME` | `268` |
| `// Key: Home
    KEY_END` | `269` |
| `// Key: End
    KEY_CAPS_LOCK` | `280` |
| `// Key: Caps lock
    KEY_SCROLL_LOCK` | `281` |
| `// Key: Scroll down
    KEY_NUM_LOCK` | `282` |
| `// Key: Num lock
    KEY_PRINT_SCREEN` | `283` |
| `// Key: Print screen
    KEY_PAUSE` | `284` |
| `// Key: Pause
    KEY_F1` | `290` |
| `// Key: F1
    KEY_F2` | `291` |
| `// Key: F2
    KEY_F3` | `292` |
| `// Key: F3
    KEY_F4` | `293` |
| `// Key: F4
    KEY_F5` | `294` |
| `// Key: F5
    KEY_F6` | `295` |
| `// Key: F6
    KEY_F7` | `296` |
| `// Key: F7
    KEY_F8` | `297` |
| `// Key: F8
    KEY_F9` | `298` |
| `// Key: F9
    KEY_F10` | `299` |
| `// Key: F10
    KEY_F11` | `300` |
| `// Key: F11
    KEY_F12` | `301` |
| `// Key: F12
    KEY_LEFT_SHIFT` | `340` |
| `// Key: Shift left
    KEY_LEFT_CONTROL` | `341` |
| `// Key: Control left
    KEY_LEFT_ALT` | `342` |
| `// Key: Alt left
    KEY_LEFT_SUPER` | `343` |
| `// Key: Super left
    KEY_RIGHT_SHIFT` | `344` |
| `// Key: Shift right
    KEY_RIGHT_CONTROL` | `345` |
| `// Key: Control right
    KEY_RIGHT_ALT` | `346` |
| `// Key: Alt right
    KEY_RIGHT_SUPER` | `347` |
| `// Key: Super right
    KEY_KB_MENU` | `348` |
| `// Key: KB menu
    // Keypad keys
    KEY_KP_0` | `320` |
| `// Key: Keypad 0
    KEY_KP_1` | `321` |
| `// Key: Keypad 1
    KEY_KP_2` | `322` |
| `// Key: Keypad 2
    KEY_KP_3` | `323` |
| `// Key: Keypad 3
    KEY_KP_4` | `324` |
| `// Key: Keypad 4
    KEY_KP_5` | `325` |
| `// Key: Keypad 5
    KEY_KP_6` | `326` |
| `// Key: Keypad 6
    KEY_KP_7` | `327` |
| `// Key: Keypad 7
    KEY_KP_8` | `328` |
| `// Key: Keypad 8
    KEY_KP_9` | `329` |
| `// Key: Keypad 9
    KEY_KP_DECIMAL` | `330` |
| `// Key: Keypad .
    KEY_KP_DIVIDE` | `331` |
| `// Key: Keypad /
    KEY_KP_MULTIPLY` | `332` |
| `// Key: Keypad *
    KEY_KP_SUBTRACT` | `333` |
| `// Key: Keypad -
    KEY_KP_ADD` | `334` |
| `// Key: Keypad +
    KEY_KP_ENTER` | `335` |
| `// Key: Keypad Enter
    KEY_KP_EQUAL` | `336` |
| `// Key: Keypad` | `108` |
| `// Key: Android back button
    KEY_MENU` | `5` |
| `// Key: Android menu button
    KEY_VOLUME_UP` | `24` |
| `// Key: Android volume up button
    KEY_VOLUME_DOWN` | `25` |

---

### MouseButton

| Name | Value |
|------|-------|
| `MOUSE_BUTTON_LEFT` | `0` |
| `// Mouse button left
    MOUSE_BUTTON_RIGHT` | `1` |
| `// Mouse button right
    MOUSE_BUTTON_MIDDLE` | `2` |
| `// Mouse button middle (pressed wheel)
    MOUSE_BUTTON_SIDE` | `3` |
| `// Mouse button side (advanced mouse device)
    MOUSE_BUTTON_EXTRA` | `4` |
| `// Mouse button extra (advanced mouse device)
    MOUSE_BUTTON_FORWARD` | `5` |
| `// Mouse button forward (advanced mouse device)
    MOUSE_BUTTON_BACK` | `6` |
| `// Mouse button back (advanced mouse device)` | `7` |

---

### MouseCursor

| Name | Value |
|------|-------|
| `MOUSE_CURSOR_DEFAULT` | `0` |
| `// Default pointer shape
    MOUSE_CURSOR_ARROW` | `1` |
| `// Arrow shape
    MOUSE_CURSOR_IBEAM` | `2` |
| `// Text writing cursor shape
    MOUSE_CURSOR_CROSSHAIR` | `3` |
| `// Cross shape
    MOUSE_CURSOR_POINTING_HAND` | `4` |
| `// Pointing hand cursor
    MOUSE_CURSOR_RESIZE_EW` | `5` |
| `// Horizontal resize/move arrow shape
    MOUSE_CURSOR_RESIZE_NS` | `6` |
| `// Vertical resize/move arrow shape
    MOUSE_CURSOR_RESIZE_NWSE` | `7` |
| `// Top-left to bottom-right diagonal resize/move arrow shape
    MOUSE_CURSOR_RESIZE_NESW` | `8` |
| `// The top-right to bottom-left diagonal resize/move arrow shape
    MOUSE_CURSOR_RESIZE_ALL` | `9` |
| `// The omnidirectional resize/move cursor shape
    MOUSE_CURSOR_NOT_ALLOWED` | `10` |

---

### GamepadButton

| Name | Value |
|------|-------|
| `GAMEPAD_BUTTON_UNKNOWN` | `0` |
| `// Unknown button` | `1` |
| `for error checking
    GAMEPAD_BUTTON_LEFT_FACE_UP` | `2` |
| `// Gamepad left DPAD up button
    GAMEPAD_BUTTON_LEFT_FACE_RIGHT` | `3` |
| `// Gamepad left DPAD right button
    GAMEPAD_BUTTON_LEFT_FACE_DOWN` | `4` |
| `// Gamepad left DPAD down button
    GAMEPAD_BUTTON_LEFT_FACE_LEFT` | `5` |
| `// Gamepad left DPAD left button
    GAMEPAD_BUTTON_RIGHT_FACE_UP` | `6` |
| `// Gamepad right button up (i.e. PS3: Triangle` | `7` |
| `Xbox: Y)
    GAMEPAD_BUTTON_RIGHT_FACE_RIGHT` | `8` |
| `// Gamepad right button right (i.e. PS3: Circle` | `9` |
| `Xbox: B)
    GAMEPAD_BUTTON_RIGHT_FACE_DOWN` | `10` |
| `// Gamepad right button down (i.e. PS3: Cross` | `11` |
| `Xbox: A)
    GAMEPAD_BUTTON_RIGHT_FACE_LEFT` | `12` |
| `// Gamepad right button left (i.e. PS3: Square` | `13` |
| `Xbox: X)
    GAMEPAD_BUTTON_LEFT_TRIGGER_1` | `14` |
| `// Gamepad top/back trigger left (first)` | `15` |
| `it could be a trailing button
    GAMEPAD_BUTTON_LEFT_TRIGGER_2` | `16` |
| `// Gamepad top/back trigger left (second)` | `17` |
| `it could be a trailing button
    GAMEPAD_BUTTON_RIGHT_TRIGGER_1` | `18` |
| `// Gamepad top/back trigger right (first)` | `19` |
| `it could be a trailing button
    GAMEPAD_BUTTON_RIGHT_TRIGGER_2` | `20` |
| `// Gamepad top/back trigger right (second)` | `21` |
| `it could be a trailing button
    GAMEPAD_BUTTON_MIDDLE_LEFT` | `22` |
| `// Gamepad center buttons` | `23` |
| `left one (i.e. PS3: Select)
    GAMEPAD_BUTTON_MIDDLE` | `24` |
| `// Gamepad center buttons` | `25` |
| `middle one (i.e. PS3: PS` | `26` |
| `Xbox: XBOX)
    GAMEPAD_BUTTON_MIDDLE_RIGHT` | `27` |
| `// Gamepad center buttons` | `28` |
| `right one (i.e. PS3: Start)
    GAMEPAD_BUTTON_LEFT_THUMB` | `29` |
| `// Gamepad joystick pressed button left
    GAMEPAD_BUTTON_RIGHT_THUMB          // Gamepad joystick pressed button right` | `30` |

---

### GamepadAxis

| Name | Value |
|------|-------|
| `GAMEPAD_AXIS_LEFT_X` | `0` |
| `// Gamepad left stick X axis
    GAMEPAD_AXIS_LEFT_Y` | `1` |
| `// Gamepad left stick Y axis
    GAMEPAD_AXIS_RIGHT_X` | `2` |
| `// Gamepad right stick X axis
    GAMEPAD_AXIS_RIGHT_Y` | `3` |
| `// Gamepad right stick Y axis
    GAMEPAD_AXIS_LEFT_TRIGGER` | `4` |
| `// Gamepad back trigger left` | `5` |
| `pressure level: [1..-1]
    GAMEPAD_AXIS_RIGHT_TRIGGER` | `5` |
| `pressure level: [1..-1]` | `7` |

---

### MaterialMapIndex

| Name | Value |
|------|-------|
| `MATERIAL_MAP_ALBEDO` | `0` |
| `// Albedo material (same as: MATERIAL_MAP_DIFFUSE)
    MATERIAL_MAP_METALNESS` | `1` |
| `// Metalness material (same as: MATERIAL_MAP_SPECULAR)
    MATERIAL_MAP_NORMAL` | `2` |
| `// Normal material
    MATERIAL_MAP_ROUGHNESS` | `3` |
| `// Roughness material
    MATERIAL_MAP_OCCLUSION` | `4` |
| `// Ambient occlusion material
    MATERIAL_MAP_EMISSION` | `5` |
| `// Emission material
    MATERIAL_MAP_HEIGHT` | `6` |
| `// Heightmap material
    MATERIAL_MAP_CUBEMAP` | `7` |
| `// Cubemap material (NOTE: Uses GL_TEXTURE_CUBE_MAP)
    MATERIAL_MAP_IRRADIANCE` | `8` |
| `// Irradiance material (NOTE: Uses GL_TEXTURE_CUBE_MAP)
    MATERIAL_MAP_PREFILTER` | `9` |
| `// Prefilter material (NOTE: Uses GL_TEXTURE_CUBE_MAP)
    MATERIAL_MAP_BRDF               // Brdf material` | `10` |

---

### ShaderLocationIndex

| Name | Value |
|------|-------|
| `SHADER_LOC_VERTEX_POSITION` | `0` |
| `// Shader location: vertex attribute: position
    SHADER_LOC_VERTEX_TEXCOORD01` | `1` |
| `// Shader location: vertex attribute: texcoord01
    SHADER_LOC_VERTEX_TEXCOORD02` | `2` |
| `// Shader location: vertex attribute: texcoord02
    SHADER_LOC_VERTEX_NORMAL` | `3` |
| `// Shader location: vertex attribute: normal
    SHADER_LOC_VERTEX_TANGENT` | `4` |
| `// Shader location: vertex attribute: tangent
    SHADER_LOC_VERTEX_COLOR` | `5` |
| `// Shader location: vertex attribute: color
    SHADER_LOC_MATRIX_MVP` | `6` |
| `// Shader location: matrix uniform: model-view-projection
    SHADER_LOC_MATRIX_VIEW` | `7` |
| `// Shader location: matrix uniform: view (camera transform)
    SHADER_LOC_MATRIX_PROJECTION` | `8` |
| `// Shader location: matrix uniform: projection
    SHADER_LOC_MATRIX_MODEL` | `9` |
| `// Shader location: matrix uniform: model (transform)
    SHADER_LOC_MATRIX_NORMAL` | `10` |
| `// Shader location: matrix uniform: normal
    SHADER_LOC_VECTOR_VIEW` | `11` |
| `// Shader location: vector uniform: view
    SHADER_LOC_COLOR_DIFFUSE` | `12` |
| `// Shader location: vector uniform: diffuse color
    SHADER_LOC_COLOR_SPECULAR` | `13` |
| `// Shader location: vector uniform: specular color
    SHADER_LOC_COLOR_AMBIENT` | `14` |
| `// Shader location: vector uniform: ambient color
    SHADER_LOC_MAP_ALBEDO` | `15` |
| `// Shader location: sampler2d texture: albedo (same as: SHADER_LOC_MAP_DIFFUSE)
    SHADER_LOC_MAP_METALNESS` | `16` |
| `// Shader location: sampler2d texture: metalness (same as: SHADER_LOC_MAP_SPECULAR)
    SHADER_LOC_MAP_NORMAL` | `17` |
| `// Shader location: sampler2d texture: normal
    SHADER_LOC_MAP_ROUGHNESS` | `18` |
| `// Shader location: sampler2d texture: roughness
    SHADER_LOC_MAP_OCCLUSION` | `19` |
| `// Shader location: sampler2d texture: occlusion
    SHADER_LOC_MAP_EMISSION` | `20` |
| `// Shader location: sampler2d texture: emission
    SHADER_LOC_MAP_HEIGHT` | `21` |
| `// Shader location: sampler2d texture: heightmap
    SHADER_LOC_MAP_CUBEMAP` | `22` |
| `// Shader location: samplerCube texture: cubemap
    SHADER_LOC_MAP_IRRADIANCE` | `23` |
| `// Shader location: samplerCube texture: irradiance
    SHADER_LOC_MAP_PREFILTER` | `24` |
| `// Shader location: samplerCube texture: prefilter
    SHADER_LOC_MAP_BRDF` | `25` |
| `// Shader location: sampler2d texture: brdf
    SHADER_LOC_VERTEX_BONEIDS` | `26` |
| `// Shader location: vertex attribute: bone indices
    SHADER_LOC_VERTEX_BONEWEIGHTS` | `27` |
| `// Shader location: vertex attribute: bone weights
    SHADER_LOC_MATRIX_BONETRANSFORMS` | `28` |
| `// Shader location: matrix attribute: bone transforms (animation)
    SHADER_LOC_VERTEX_INSTANCETRANSFORM // Shader location: vertex attribute: instance transforms` | `29` |

---

### ShaderUniformDataType

| Name | Value |
|------|-------|
| `SHADER_UNIFORM_FLOAT` | `0` |
| `// Shader uniform type: float
    SHADER_UNIFORM_VEC2` | `1` |
| `// Shader uniform type: vec2 (2 float)
    SHADER_UNIFORM_VEC3` | `2` |
| `// Shader uniform type: vec3 (3 float)
    SHADER_UNIFORM_VEC4` | `3` |
| `// Shader uniform type: vec4 (4 float)
    SHADER_UNIFORM_INT` | `4` |
| `// Shader uniform type: int
    SHADER_UNIFORM_IVEC2` | `5` |
| `// Shader uniform type: ivec2 (2 int)
    SHADER_UNIFORM_IVEC3` | `6` |
| `// Shader uniform type: ivec3 (3 int)
    SHADER_UNIFORM_IVEC4` | `7` |
| `// Shader uniform type: ivec4 (4 int)
    SHADER_UNIFORM_UINT` | `8` |
| `// Shader uniform type: unsigned int
    SHADER_UNIFORM_UIVEC2` | `9` |
| `// Shader uniform type: uivec2 (2 unsigned int)
    SHADER_UNIFORM_UIVEC3` | `10` |
| `// Shader uniform type: uivec3 (3 unsigned int)
    SHADER_UNIFORM_UIVEC4` | `11` |
| `// Shader uniform type: uivec4 (4 unsigned int)
    SHADER_UNIFORM_SAMPLER2D        // Shader uniform type: sampler2d` | `12` |

---

### ShaderAttributeDataType

| Name | Value |
|------|-------|
| `SHADER_ATTRIB_FLOAT` | `0` |
| `// Shader attribute type: float
    SHADER_ATTRIB_VEC2` | `1` |
| `// Shader attribute type: vec2 (2 float)
    SHADER_ATTRIB_VEC3` | `2` |
| `// Shader attribute type: vec3 (3 float)
    SHADER_ATTRIB_VEC4              // Shader attribute type: vec4 (4 float)` | `3` |

---

### PixelFormat

| Name | Value |
|------|-------|
| `PIXELFORMAT_UNCOMPRESSED_GRAYSCALE` | `1` |
| `// 8 bit per pixel (no alpha)
    PIXELFORMAT_UNCOMPRESSED_GRAY_ALPHA` | `1` |
| `// 8*2 bpp (2 channels)
    PIXELFORMAT_UNCOMPRESSED_R5G6B5` | `2` |
| `// 16 bpp
    PIXELFORMAT_UNCOMPRESSED_R8G8B8` | `3` |
| `// 24 bpp
    PIXELFORMAT_UNCOMPRESSED_R5G5B5A1` | `4` |
| `// 16 bpp (1 bit alpha)
    PIXELFORMAT_UNCOMPRESSED_R4G4B4A4` | `5` |
| `// 16 bpp (4 bit alpha)
    PIXELFORMAT_UNCOMPRESSED_R8G8B8A8` | `6` |
| `// 32 bpp
    PIXELFORMAT_UNCOMPRESSED_R32` | `7` |
| `// 32 bpp (1 channel - float)
    PIXELFORMAT_UNCOMPRESSED_R32G32B32` | `8` |
| `// 32*3 bpp (3 channels - float)
    PIXELFORMAT_UNCOMPRESSED_R32G32B32A32` | `9` |
| `// 32*4 bpp (4 channels - float)
    PIXELFORMAT_UNCOMPRESSED_R16` | `10` |
| `// 16 bpp (1 channel - half float)
    PIXELFORMAT_UNCOMPRESSED_R16G16B16` | `11` |
| `// 16*3 bpp (3 channels - half float)
    PIXELFORMAT_UNCOMPRESSED_R16G16B16A16` | `12` |
| `// 16*4 bpp (4 channels - half float)
    PIXELFORMAT_COMPRESSED_DXT1_RGB` | `13` |
| `// 4 bpp (no alpha)
    PIXELFORMAT_COMPRESSED_DXT1_RGBA` | `14` |
| `// 4 bpp (1 bit alpha)
    PIXELFORMAT_COMPRESSED_DXT3_RGBA` | `15` |
| `// 8 bpp
    PIXELFORMAT_COMPRESSED_DXT5_RGBA` | `16` |
| `// 8 bpp
    PIXELFORMAT_COMPRESSED_ETC1_RGB` | `17` |
| `// 4 bpp
    PIXELFORMAT_COMPRESSED_ETC2_RGB` | `18` |
| `// 4 bpp
    PIXELFORMAT_COMPRESSED_ETC2_EAC_RGBA` | `19` |
| `// 8 bpp
    PIXELFORMAT_COMPRESSED_PVRT_RGB` | `20` |
| `// 4 bpp
    PIXELFORMAT_COMPRESSED_PVRT_RGBA` | `21` |
| `// 4 bpp
    PIXELFORMAT_COMPRESSED_ASTC_4x4_RGBA` | `22` |
| `// 8 bpp
    PIXELFORMAT_COMPRESSED_ASTC_8x8_RGBA    // 2 bpp` | `23` |

---

### TextureFilter

| Name | Value |
|------|-------|
| `TEXTURE_FILTER_POINT` | `0` |
| `// No filter` | `1` |
| `pixel approximation
    TEXTURE_FILTER_BILINEAR` | `2` |
| `// Linear filtering
    TEXTURE_FILTER_TRILINEAR` | `3` |
| `// Trilinear filtering (linear with mipmaps)
    TEXTURE_FILTER_ANISOTROPIC_4X` | `4` |
| `// Anisotropic filtering 4x
    TEXTURE_FILTER_ANISOTROPIC_8X` | `5` |
| `// Anisotropic filtering 8x
    TEXTURE_FILTER_ANISOTROPIC_16X` | `6` |
| `// Anisotropic filtering 16x` | `7` |

---

### TextureWrap

| Name | Value |
|------|-------|
| `TEXTURE_WRAP_REPEAT` | `0` |
| `// Repeats texture in tiled mode
    TEXTURE_WRAP_CLAMP` | `1` |
| `// Clamps texture to edge pixel in tiled mode
    TEXTURE_WRAP_MIRROR_REPEAT` | `2` |
| `// Mirrors and repeats the texture in tiled mode
    TEXTURE_WRAP_MIRROR_CLAMP               // Mirrors and clamps to border the texture in tiled mode` | `3` |

---

### CubemapLayout

| Name | Value |
|------|-------|
| `CUBEMAP_LAYOUT_AUTO_DETECT` | `0` |
| `// Automatically detect layout type
    CUBEMAP_LAYOUT_LINE_VERTICAL` | `1` |
| `// Layout is defined by a vertical line with faces
    CUBEMAP_LAYOUT_LINE_HORIZONTAL` | `2` |
| `// Layout is defined by a horizontal line with faces
    CUBEMAP_LAYOUT_CROSS_THREE_BY_FOUR` | `3` |
| `// Layout is defined by a 3x4 cross with cubemap faces
    CUBEMAP_LAYOUT_CROSS_FOUR_BY_THREE     // Layout is defined by a 4x3 cross with cubemap faces` | `4` |

---

### FontType

| Name | Value |
|------|-------|
| `FONT_DEFAULT` | `0` |
| `// Default font generation` | `1` |
| `anti-aliased
    FONT_BITMAP` | `2` |
| `// Bitmap font generation` | `3` |
| `no anti-aliasing
    FONT_SDF                        // SDF font generation` | `4` |
| `requires external shader` | `5` |

---

### BlendMode

| Name | Value |
|------|-------|
| `BLEND_ALPHA` | `0` |
| `// Blend textures considering alpha (default)
    BLEND_ADDITIVE` | `1` |
| `// Blend textures adding colors
    BLEND_MULTIPLIED` | `2` |
| `// Blend textures multiplying colors
    BLEND_ADD_COLORS` | `3` |
| `// Blend textures adding colors (alternative)
    BLEND_SUBTRACT_COLORS` | `4` |
| `// Blend textures subtracting colors (alternative)
    BLEND_ALPHA_PREMULTIPLY` | `5` |
| `// Blend premultiplied textures considering alpha
    BLEND_CUSTOM` | `6` |
| `// Blend textures using custom src/dst factors (use rlSetBlendFactors())
    BLEND_CUSTOM_SEPARATE           // Blend textures using custom rgb/alpha separate src/dst factors (use rlSetBlendFactorsSeparate())` | `7` |

---

### Gesture

| Name | Value |
|------|-------|
| `GESTURE_NONE` | `0` |
| `// No gesture
    GESTURE_TAP` | `1` |
| `// Tap gesture
    GESTURE_DOUBLETAP` | `2` |
| `// Double tap gesture
    GESTURE_HOLD` | `4` |
| `// Hold gesture
    GESTURE_DRAG` | `8` |
| `// Drag gesture
    GESTURE_SWIPE_RIGHT` | `16` |
| `// Swipe right gesture
    GESTURE_SWIPE_LEFT` | `32` |
| `// Swipe left gesture
    GESTURE_SWIPE_UP` | `64` |
| `// Swipe up gesture
    GESTURE_SWIPE_DOWN` | `128` |
| `// Swipe down gesture
    GESTURE_PINCH_IN` | `256` |
| `// Pinch in gesture
    GESTURE_PINCH_OUT` | `512` |

---

### CameraMode

| Name | Value |
|------|-------|
| `CAMERA_CUSTOM` | `0` |
| `// Camera custom` | `1` |
| `controlled by user (UpdateCamera() does nothing)
    CAMERA_FREE` | `2` |
| `// Camera free mode
    CAMERA_ORBITAL` | `3` |
| `// Camera orbital` | `4` |
| `around target` | `5` |
| `zoom supported
    CAMERA_FIRST_PERSON` | `6` |
| `// Camera first person
    CAMERA_THIRD_PERSON             // Camera third person` | `7` |

---

### CameraProjection

| Name | Value |
|------|-------|
| `CAMERA_PERSPECTIVE` | `0` |
| `// Perspective projection
    CAMERA_ORTHOGRAPHIC             // Orthographic projection` | `1` |

---

### NPatchLayout

| Name | Value |
|------|-------|
| `NPATCH_NINE_PATCH` | `0` |
| `// Npatch layout: 3x3 tiles
    NPATCH_THREE_PATCH_VERTICAL` | `1` |
| `// Npatch layout: 1x3 tiles
    NPATCH_THREE_PATCH_HORIZONTAL   // Npatch layout: 3x1 tiles` | `2` |

---

## Macros

### RAYLIB_H

```c
#define RAYLIB_H 
```

---

### RAYLIB_VERSION_MAJOR

```c
#define RAYLIB_VERSION_MAJOR 
```

---

### RAYLIB_VERSION_MINOR

```c
#define RAYLIB_VERSION_MINOR 
```

---

### RAYLIB_VERSION_PATCH

```c
#define RAYLIB_VERSION_PATCH 
```

---

### RAYLIB_VERSION

```c
#define RAYLIB_VERSION 
```

---

### RL_COLOR_TYPE

```c
#define RL_COLOR_TYPE 
```

---

### RL_VECTOR2_TYPE

```c
#define RL_VECTOR2_TYPE 
```

---

### RL_VECTOR4_TYPE

```c
#define RL_VECTOR4_TYPE 
```

---

### RL_MATRIX_TYPE

```c
#define RL_MATRIX_TYPE 
```

---

### LIGHTGRAY

```c
#define LIGHTGRAY 
```

---

### GRAY

```c
#define GRAY 
```

---

### DARKGRAY

```c
#define DARKGRAY 
```

---

### YELLOW

```c
#define YELLOW 
```

---

### GOLD

```c
#define GOLD 
```

---

### ORANGE

```c
#define ORANGE 
```

---

### PINK

```c
#define PINK 
```

---

### RED

```c
#define RED 
```

---

### MAROON

```c
#define MAROON 
```

---

### GREEN

```c
#define GREEN 
```

---

### LIME

```c
#define LIME 
```

---

### DARKGREEN

```c
#define DARKGREEN 
```

---

### SKYBLUE

```c
#define SKYBLUE 
```

---

### BLUE

```c
#define BLUE 
```

---

### DARKBLUE

```c
#define DARKBLUE 
```

---

### PURPLE

```c
#define PURPLE 
```

---

### VIOLET

```c
#define VIOLET 
```

---

### DARKPURPLE

```c
#define DARKPURPLE 
```

---

### BEIGE

```c
#define BEIGE 
```

---

### BROWN

```c
#define BROWN 
```

---

### DARKBROWN

```c
#define DARKBROWN 
```

---

### WHITE

```c
#define WHITE 
```

---

### BLACK

```c
#define BLACK 
```

---

### BLANK

```c
#define BLANK 
```

---

### MAGENTA

```c
#define MAGENTA 
```

---

### RAYWHITE

```c
#define RAYWHITE 
```

---

### MOUSE_LEFT_BUTTON

```c
#define MOUSE_LEFT_BUTTON 
```

---

### MOUSE_RIGHT_BUTTON

```c
#define MOUSE_RIGHT_BUTTON 
```

---

### MOUSE_MIDDLE_BUTTON

```c
#define MOUSE_MIDDLE_BUTTON 
```

---

### MATERIAL_MAP_DIFFUSE

```c
#define MATERIAL_MAP_DIFFUSE 
```

---

### MATERIAL_MAP_SPECULAR

```c
#define MATERIAL_MAP_SPECULAR 
```

---

### SHADER_LOC_MAP_DIFFUSE

```c
#define SHADER_LOC_MAP_DIFFUSE 
```

---

### SHADER_LOC_MAP_SPECULAR

```c
#define SHADER_LOC_MAP_SPECULAR 
```

---

### GetMouseRay

```c
#define GetMouseRay 
```

---

