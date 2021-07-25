# WTHE
This module attempts to enhance contrast of a given image or video by employing a method called weighted thresholded histogram equalization (WTHE). This method seeks to improve on preformance of the conventional histogram equalization method by adding controllable parameters to it. By weighting and thresholding the PMF of the image before performing histogram equalization, two parameters are introduced that can be changed manually, but by experimenting on a variety of images, optimal values for both parameters are calculated (r = 0.5, v = 0.5).


You can access the article that came up with this method [here](https://www.researchgate.net/publication/3183125_Ward_RK_Fast_ImageVideo_Contrast_Enhancement_Based_on_Weighted_Thresholded_Histogram_Equalization_IEEE_Trans_Consumer_Electronics_532_757-764). 


## Showcase
* A sample image and the sample image after enhancement  
![cloudy-day jpg](https://github.com/Mamdasn/wthe/blob/main/showcase/assets/cloudy-day-400.jpg "cloudy-day jpg") ![Enhanced-cloudy-day jpg](https://github.com/Mamdasn/wthe/blob/main/showcase/output/Enhanced-cloudy-day-400.jpg "Enhanced-cloudy-day jpg") 

## Future works
* disassemble wthe.go into seperate categories of functions
* add more comments