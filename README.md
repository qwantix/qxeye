# QxEye

QxEye is a KISS-coded project to ensure that the threat is there before reporting the incident.
It is not and there is no remote access system. **QxEye is a stand-alone vision system**.

It is designed to run on tiny hardware like RaspberryPI or OrangePI, but requires at least an ARMv7 to take advantage of hardware optimizations.


## How It Works

QxEye reads the camera's video stream. It will compare the image with a stack of previous images to isolate the movement zones and then analyze each zones with a pre-trained dnn to detect objects with precision, such as a person.

Finally, you will have the choice to save a capture, tweet direct message, call an url or execute a script when a match is positive.

## Configuration

Configuration is splitted in 4 sections :
- `cameras` : list of cameras streams
- `matchers`: list of matchers used to provide detection
- `triggers`: list of triggers 
- `services`: list of services configuration

### Cameras

`cameras` section is an array of object like this

```js
 "cameras": [
    {
      "name": String, // Name of camera
      "enabled": Boolean, // Camera enabled or not
      "endpoint": String, // Endpoint of camera
      "persistence": Int, // Image persistence, see below
      "matcher": String, // Matcher name 
      "zones": [ // Array of zone
        {
          "ignore": Bool, // Optional, indicate if zone is opaque or confidential
          "color": String, // Color (RRGGBB) of non ignored zone
          "mask": [ 
            // Mask definition, see below
          ]
        }
      ]
    }
  ]
```

#### Image Persistence
Image persitence is used to isolate motion in image.

When a new image is read, it merges with the previous images with a transparency factor to create a new image representing the average of movements such as wind in trees or grass...

`persistence` is the thousandth of an opacity applied to the image before merging into the stack.

Therefore, higher is the value, more sensitivity will increase but may cause CPU overload.

#### Zones
Zones have 2 utilities, firstly ignore sensitive or uninteresting places, such as the neighbor's garden. Secondly, define zones in which detections can have a different meaning.

`mask` is an array of string, each row represent horizontal slice on image, each character, represent au vertical slice of row.

A blank " " char represent a empty space, and char "#" represent a mask

This mask below represent a mask than slice image in 2 row of 2 columns and hide the first quarter
```js
[
    "# ",
    "  "
]
```

But each column is relative to its line, so you can write that to do the same thing.
```js
[
    "# ",
    ""
]
```

For example, more complicated :
```js
[
    "## ##",
    ""
    ""
    " #"
]
```
Image will be sliced into 4 lines, first line will be slices in 5 cols and last in 2 col

### Matchers

`matchers` is an array of object like this:

```js
"matchers" : [
    {
      "name": String, // Name of matcher
      "type": String, // Type of matcher, set to "dnn"
      "params": { 
        // Object of params
      }
    }
]
```

### Triggers

`triggers` is an array of object like this:

```js
"triggers": [
    {
      "on": String,
      "zones": [ String ], // Optionnal
      "confidence": Float, // confidence level 
      "service": String, // Service used
      "delay": Int, // Time delay in seconds during which the trigger is deactivated after being called
      "params": {
          // Object of trigger params
      }
    }
]
```

### Services

`services` is an object of object like this:

```js
"services": [
   "serviceName": { 
       // Params
   }
]
```
#### file
Capture image and save it to a file

```js
{
    "dir": "./captures", // Dir to store file
    "traces": "matches,ignoredZonesOnly" // Comma separated trace flag
}
```
**traces** add some traces to the capture  and can take value:
* `matches` print rect of item matches
* `ignoredZonesOnly` draw an black rect on the ignored zones
* `zones` print zones


#### twitter
```js
{
    "consumerKey": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "consumerSecret": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "accessToken": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "accessTokenSecret": "xxxxxxxxxxxxxxxxxxxxxxxx"
}
```

#### url
Not implemented yet

#### script
Not implemented yet

---

## Requirements

* Go 1.9
* OpenCV 3.3


## Roadmap & Notes

* Implement url and script trigger
* Support YOLO2 model
* Support haarcascade ? useful?
* Record a video sample instead of capture
* Pretrain model optimized for video surveillance
* Auto adjust trained model
* Face detection ?
* Detect dangerous, hesitant or suspicious trajectories



