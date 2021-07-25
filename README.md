# WTHE
This module attempts to enhance contrast of a given image or video by employing a method called weighted thresholded histogram equalization (WTHE). This method seeks to improve on preformance of the conventional histogram equalization method by adding controllable parameters to it. By weighting and thresholding the PMF of the image before performing histogram equalization, two parameters are introduced that can be changed manually, but by experimenting on a variety of images, optimal values for both parameters are calculated (r = 0.5, v = 0.5).


You can access the article that came up with this method [here](https://www.researchgate.net/publication/3183125_Ward_RK_Fast_ImageVideo_Contrast_Enhancement_Based_on_Weighted_Thresholded_Histogram_Equalization_IEEE_Trans_Consumer_Electronics_532_757-764). 


## Showcase
* A sample image and the sample image after enhancement  
![cloudy-day.jpg jpg](https://raw.githubusercontent.com/Mamdasn/wthe/main/showcase/assets/cloudy-day-small.jpg "cloudy-day.jpg jpg") ![Enhanced-cloudy-day.jpg jpg](https://raw.githubusercontent.com/Mamdasn/wthe/main/showcase/output/Enhanced-cloudy-day-small.jpg "Enhanced-cloudy-day.jpg jpg") 
