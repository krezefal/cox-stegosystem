# Cox Algorithm

## Description

This algorithm is robust to many signal processing operations. Detection of an embedded digital watermark in it is performed using the original image. The embedded data is a sequence of real numbers with zero mean and unit variance. Several ACs are used to embed information.DCT coefficients of the entire image with the highest energy.

The first option can be used in the case when the energy of the CVZ is comparable to the energy of the modified coefficient. Otherwise, either the digital water center will be non-robust, or the distortion is too large. Therefore, it is possible to embed information in this way only with an insignificant range of variation in the values ​​of the energy of the coefficients.When a digital watermark is detected, the reverse operations are performed: the DCTs of the original and modified images are calculated, the differences between the corresponding coefficients of the largest value are found.