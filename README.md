# QxEye

QxEye is a KISS-coded project to ensure that the threat is there before reporting the incident.
It is not and there is no remote access system. **QxEye is a stand-alone vision system**.

It is designed to run on tiny hardware like RaspberryPI or OrangePI, but requires at least an ARMv7 to take advantage of hardware optimizations.


## How It Works

QxEye reads the camera's video feed. It will compare the image with a stack of previous images to isolate the movement zones and then analyze each zones with a pre-trained dnn to detect objects with precision, such as a person.

Finally, you will have the choice to save a capture, tweet direct message, call an url or execute a script when a match is positive.

---

# Requirements

* Go 1.9
* OpenCV 3.3

# Roadmap

* Implement url and script trigger
* Support YOLO2 model
* Support haarcascade ? useful?
* Record a video sample instead of capture
* Pretrain model optimized for video surveillance
* Auto adjust trained model
* Face detection ?
* Detect dangerous, hesitant or suspicious trajectories



