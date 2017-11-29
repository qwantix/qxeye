#ifndef CONFIG_H
#define CONFIG_H

// Simple config class
// Reimplemented from https://wiki.calculquebec.ca/w/C%2B%2B_:_fichier_de_configuration/en

#include <map>
#include <string>
#include <istream>

using namespace std;

class Config
{
public:
  Config();
  bool loadFile(string filename);
  bool loadString(string config);
  bool loadStream(istream &stream);

  bool contains(const string &key) const;
  bool get(const string &key, string &value) const;
  bool get(const string &key, int &value) const;
  bool get(const string &key, float &value) const;
  bool get(const string &key, bool &value) const;

protected:
  string trim(const string str) const;

private:
  map<string, string> data;
};

#endif //CONFIG_H
