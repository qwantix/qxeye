#include "util.hpp"

#include <chrono>
#include <iostream>
#include <string>
#include <map>

using namespace std;

std::map<std::string, std::chrono::high_resolution_clock::time_point> timers;

void startTimer(std::string name)
{
  timers[name] = std::chrono::high_resolution_clock::now();
}

void endTimer(std::string name)
{
  std::chrono::high_resolution_clock::time_point t2 = std::chrono::high_resolution_clock::now();
  std::chrono::duration<double, std::milli> time_span = t2 - timers[name];
  cout << name << " " << time_span.count() << "ms" << endl;
}