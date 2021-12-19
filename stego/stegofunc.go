package pkg

import (
	"image"
)

func EmbedMessage(img image.Image, message byte[], alpha float32) (image.Image, int) {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	size_t
	message_len = message.size();

	if (message_len > (data.size()/8)*(data[0].size()/8)) {

	std::string
		error_mes = "Invalid arguments! Image capacity: " +
			std::to_string((data.size()/8)*(data[0].size()/8)) +
			", container capacity" + std::to_string(message_len);

		throw
	std::length_error(error_mes);
	}

std::vector < std::vector < long
	double >> res(data.size(), std::vector < long
	double > (data[0].size()));

	size_t
	counter = 0;

	for (std::size_t
	h = 0;
	h < data.size();
	h += N){
		for (std::size_t
		w = 0;
		w < data[0].size();
		w += N){
			int
			max_i = 0;
			int
			max_j = 1;
			for (int
			k = 0;
			k < N;
			++k){
				long
				double
				Ck = (k == 0 ? (1.0 / N) : (2.0 / N));
				for (int
				l = 0;
				l < N;
				++l){
					long
					double
					Cl = (l == 0 ? (1.0 / N) : (2.0 / N));
					long
					double
					sum = 0;
					for (int
					i = 0;
					i < N;
					++i){
						long
						double
						cos_x = cosl(((2*i + 1) * M_PI * k) / (2 * N));
						for (int
						j = 0;
						j < N;
						++j){
							sum += data[h+i][w+j] * cos_x * cosl(((2*j+1)*M_PI*l)/(2*N));
						}
					}
					long
					double
					val = sqrtl(Ck) * sqrtl(Cl) * sum;
					res[h+k][w+l] = val;
					if (counter != message_len) {
						if (k == 0 && l == 0) {
							continue;
						}
						if (std::abs(val) > std::abs(res[h+max_i][w+max_j])) {
							max_i = k;
							max_j = l;
						}
					}
				}
			}
			if (counter != message_len) {
				int
				si = message[counter] == 0 ? 1 : -1;
				res[h+max_i][w+max_j] += alpha * si;
				++counter;
			}
		}
	}

	auto
	ifdct = iFDCT(res, N);
	return ifdct;
}