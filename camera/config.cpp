#include <map>
#include <string>
#include <istream>
#include <fstream>
#include <sstream>
#include <iostream>
#include <cstdlib>
#include "config.hpp"

using namespace std;

Config::Config()
{
}

bool Config::loadFile(string filename)
{
    ifstream file(filename.c_str());
    return this->loadStream(file);
}
bool Config::loadString(string content)
{
    istringstream is(content);
    return this->loadStream(is);
}
bool Config::loadStream(istream &stream)
{
    if (!stream.good())
    {
        cerr << "Cannot read configuration" << endl;
        return false;
    }

    while (stream.good() && !stream.eof())
    {
        string line;
        getline(stream, line);

        // filter out comments
        if (!line.empty())
        {
            int pos = line.find('#');

            if (pos != string::npos)
            {
                line = line.substr(0, pos);
            }
        }

        // split line into key and value
        if (!line.empty())
        {
            int pos = line.find('=');

            if (pos != string::npos)
            {
                string key = this->trim(line.substr(0, pos));
                string value = this->trim(line.substr(pos + 1));
                if (!key.empty() && !value.empty())
                {
                    this->data[key] = value;
                }
            }
        }
    }

    return true;
}

bool Config::contains(const string &key) const
{
    return this->data.find(key) != this->data.end();
}

bool Config::get(const string &key, string &value) const
{
    map<string, string>::const_iterator iter = this->data.find(key);

    if (iter != this->data.end())
    {
        value = iter->second;
        return true;
    }
    return false;
}

bool Config::get(const string &key, int &value) const
{
    string str;
    if (this->get(key, str))
    {
        value = atoi(str.c_str());
        return true;
    }
    return false;
}

bool Config::get(const string &key, float &value) const
{
    string str;
    if (this->get(key, str))
    {
        value = atof(str.c_str());
        return true;
    }
    return false;
}

bool Config::get(const string &key, bool &value) const
{
    string str;
    if (this->get(key, str))
    {
        value = str == "true" || str == "True" || str == "1";
        return true;
    }
    return false;
}

string Config::trim(const string str) const
{
    int first = str.find_first_not_of(" \t");

    if (first != string::npos)
    {
        int last = str.find_last_not_of(" \t");
        return str.substr(first, last - first + 1);
    }
    return "";
}
