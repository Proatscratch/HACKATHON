#define SDL_MAIN_HANDLED
#include <SDL2/SDL.h>
#include <iostream>
#include <stdio.h>
#define IMGUI_DEFINE_MATH_OPERATORS
#include "imgui.h"
#include "imgui_impl_sdl2.h"
#include "imgui_impl_sdlrenderer2.h"
#include <glm/glm.hpp>
#include "imgui_internal.h"
//gui/imgui.cpp gui/imgui_impl_sdl2.cpp gui/imgui_impl_sdlrenderer2.cpp gui/imgui_internal.cpp
#define FILEDATASIZES 1000
#define TEXTSIZE 10000000
#define LETTERLIM 40

char FILEHANDLE[FILEDATASIZES];
char RawText[TEXTSIZE];
#define text_buffer RawText
bool g_insert_request = false;
char g_char_to_insert[LETTERLIM] = {0};

#define STYLESHEET_RELWINDOW_W -300
#define STYLESHEET_RELWINDOW_H 800
// 1. Define the callback function
int TextEditCallback(ImGuiInputTextCallbackData* data) {
    if (g_insert_request) {
        // data->CursorPos gives the exact index where the text cursor is currently blinking
        data->InsertChars(data->CursorPos, g_char_to_insert);
        g_insert_request = false; // Reset request
    }
    return 0;
}
void RenderUTF8TypewriterArea(ImVec2& sz, int &typing) {
    ImGui::PushStyleColor(ImGuiCol_FrameBg, ImVec4(255, 255, 255, 255));
    ImGui::PushStyleColor(ImGuiCol_FrameBgHovered, ImVec4(0, 0, 0, 0));
    ImGui::PushStyleColor(ImGuiCol_FrameBgActive, ImVec4(0, 0, 0, 0));
    ImGui::PushStyleColor(ImGuiCol_Text, ImVec4(0, 0, 0, 0));
    
    sz.x += STYLESHEET_RELWINDOW_W;
    sz.y += STYLESHEET_RELWINDOW_H;
    typing = ImGui::InputTextMultiline("##hidden_input", text_buffer, IM_ARRAYSIZE(text_buffer), sz, ImGuiInputTextFlags_CallbackAlways, TextEditCallback);
    
    sz.x -= STYLESHEET_RELWINDOW_W;
    sz.y -= STYLESHEET_RELWINDOW_H;
    ImGui::PopStyleColor(4);
    
    ImVec2 start_pos = ImGui::GetItemRectMin();
    
    ImDrawList* draw_list = ImGui::GetWindowDrawList();
    ImFont* font = ImGui::GetFont();
    float font_size = ImGui::GetFontSize();
    
    ImVec2 current_pos = start_pos;
    float time = (float)ImGui::GetTime();
    
    const char* ptr = text_buffer;
    int char_index = 0;
    
    while (*ptr != '\0') {
        unsigned int c = 0;
        
        int bytes = ImTextCharFromUtf8(&c, ptr, nullptr);
        if (*ptr == '\n') {
            current_pos.x = start_pos.x;
            
            ImVec2 char_size = ImGui::CalcTextSize("A", "A" + 1);
            current_pos.y += char_size.y;
        }
        if (bytes <= 0) {
            ptr++;
            continue; 
        }
        
        char temp_str[5] = {0};
        for (int b = 0; b < bytes && ptr[b] != '\0'; b++) {
            temp_str[b] = ptr[b];
        }
        
        ImU32 col = ImColor::HSV(0, 0, 0.0f);
        
        draw_list->AddText(font, font_size, current_pos, col, temp_str);
        
        ImVec2 char_size = ImGui::CalcTextSize(temp_str, temp_str + bytes);
        current_pos.x += char_size.x;
        
        ptr += bytes;
        char_index++;
    }
}
void insert(const char *C) {
    
        strcpy_s(g_char_to_insert,C);
        g_insert_request = true;
}
#include <stdlib.h>
float DELETE_G;
extern "C" {
float GetDeleteVal() {
    return DELETE_G;
}
void GOFUNCTION();
int fn() {
    if (SDL_Init(SDL_INIT_VIDEO) != 0) {
        printf("SDL_Init Error: %s\n", SDL_GetError());
        return 1;
    }
    
    char FILENAME[FILEDATASIZES];
    SDL_Window* window = SDL_CreateWindow(
        "SDL2 Renderer Example",
        SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED,
        800, 600,
        SDL_WINDOW_SHOWN | SDL_WINDOW_FULLSCREEN_DESKTOP
    );

    if (!window) {
        printf("SDL_CreateWindow Error: %s\n", SDL_GetError());
        SDL_Quit();
        return 1;
    }

    IMGUI_CHECKVERSION();
    ImGui::CreateContext();
    ImGuiIO& io{ImGui::GetIO()};
    
    io.ConfigFlags |= ImGuiConfigFlags_NavEnableKeyboard;
    io.ConfigFlags |= ImGuiConfigFlags_ViewportsEnable;

    SDL_Renderer* renderer = SDL_CreateRenderer(
        window, 
        -1, 
        SDL_RENDERER_ACCELERATED | SDL_RENDERER_PRESENTVSYNC
    );

    ImGui_ImplSDL2_InitForSDLRenderer(window, renderer);
    ImGui_ImplSDLRenderer2_Init(renderer);

    if (!renderer) {
        printf("SDL_CreateRenderer Error: %s\n", SDL_GetError());
        SDL_DestroyWindow(window);
        SDL_Quit();
        return 1;
    }

    SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255);
    
    ImFont* custom_font = io.Fonts->AddFontFromFileTTF("c:/Windows/Fonts/segoeui.ttf", 32.0f);
    ImFontConfig config;
    config.MergeMode = true;
    config.PixelSnapH = true;
    
    static const ImWchar emoji_ranges[] = { 0x0020, 0x00FF, 0x1F300, 0x1FAFF, 0 }; 
    io.Fonts->AddFontFromFileTTF("c:/Windows/Fonts/seguiemj.ttf", 32.0f, &config, emoji_ranges);

    io.Fonts->Build();
    int fontSize;
    float TIMER;
    float deletes = 0;
    float timer = 0.0;
    int ticks = 0;
    float time = 0;
    int CLICKS = 0;
    int TAB = 0;
    int LPM = 1;
    int QPM = 1;
    srand(SDL_GetPerformanceCounter());
    sprintf(FILENAME, "Essay%u.txt", rand());
    bool open = 1;
    int typing = 0;
    while (1) {
        int tticks = ticks;
        ticks = SDL_GetTicks();
        float dt = (ticks-tticks)/1000.0;
        if (open) {
            if (fmod(time, 0.05) > fmod(time + dt, 0.05)) {
                deletes+=CLICKS/2;
                CLICKS = 0;
            }
            if (fmod(time, 10) > fmod(time + dt,10)) {
                if (LPM/QPM < 0.5) deletes+=time * 10;
                QPM = LPM;
                LPM = 1;
            }
            if (fmod(time, 30) > fmod(time + dt,30)) {
                std::cout << deletes/time << std::endl;
            }
            time += dt;
        }
        SDL_Event event;
        while (SDL_PollEvent(&event)) {
            
            if (event.type == SDL_KEYDOWN && !event.key.repeat && open && typing) {
                LPM++;
                if (event.key.keysym.sym == SDLK_BACKSPACE || event.key.keysym.sym == SDLK_DELETE) {
                    deletes++;
                }
                
                if (event.key.keysym.sym == SDLK_TAB) {
                    insert("\t");
                }

            }
            if (event.type == SDL_KEYDOWN && typing) {
                
                CLICKS+=1;
            }
            ImGui_ImplSDL2_ProcessEvent(&event);
            if (event.type == SDL_QUIT) {
                SDL_Quit();
            }
        }
        
        ImGui_ImplSDL2_NewFrame();
        
        ImGui::NewFrame();
        ImVec2 sz;
        ImVec2 sub = sz;
        sub.y -= 10;
        ImVec2 zero = {0, 0};
        glm::ivec2 p;
        SDL_GetRendererOutputSize(renderer, &p.x, &p.y);
        sz.x = p.x;
        sz.y = p.y;
        ImGui::SetNextWindowSize(sz);
        ImGui::SetNextWindowPos(zero);
        ImGui::Begin("TextEditor", &open, ImGuiWindowFlags_MenuBar);
        if (ImGui::BeginMenuBar()) {
            
            // --- FILE MENU ---
            if (ImGui::BeginMenu("File")) {
                if (ImGui::MenuItem("New", "Ctrl+N")) {
                    FILE *file = fopen(FILENAME, "w");
                    fprintf(file, RawText);
                    fclose(file);
                    memset(FILENAME, 0, sizeof(FILENAME));
                    sprintf(FILENAME, "Essay%u.txt", rand());
                    memset(RawText, 0, sizeof(RawText));
                    memcpy(FILEHANDLE, FILENAME, sizeof(FILENAME));
                }
                if (ImGui::MenuItem("Open", "Ctrl+O")) {
                    FILE *file = fopen(FILEHANDLE, "w");
                    fprintf(file, RawText);
                    fclose(file);
                    
                    file = fopen(FILENAME, "r");
                    fread( RawText, 1, sizeof(RawText), file);
                    fclose(file);

                    memcpy(FILEHANDLE, FILENAME, sizeof(FILENAME));
                }
                if (ImGui::MenuItem("Save", "Ctrl+S")) {
                    
                    FILE *file = fopen(FILENAME, "w");
                    fprintf(file, RawText);
                    memcpy(FILEHANDLE, FILENAME, sizeof(FILEHANDLE));
                    fclose(file);
                }
                if (ImGui::MenuItem("Pause")) {
                    open = not open;
                }
                
                
                
                ImGui::EndMenu(); // Always match BeginMenu with EndMenu
            }
            
            ImGui::EndMenuBar(); // Always match BeginMenu with EndMenu
        }
        ImGui::InputText("Filename", FILENAME, sizeof(FILENAME));
        ImGui::PushItemWidth(100); 
        ImGui::Text("%f DELETES PER SECOND [FILENAME ASSIGNED TO FILE: %s]", deletes/time, FILEHANDLE);
        ImGui::PopItemWidth();
        if (open)
            RenderUTF8TypewriterArea(sz, typing);
        ImGui::End();
        ImGui::Render();
        DELETE_G = deletes/time;
        SDL_RenderClear(renderer);
        
        ImGui_ImplSDLRenderer2_RenderDrawData(ImGui::GetDrawData(), renderer);
        SDL_RenderPresent(renderer);
    }

    SDL_DestroyRenderer(renderer);
    SDL_DestroyWindow(window);
    SDL_Quit();
    return 0;
}}