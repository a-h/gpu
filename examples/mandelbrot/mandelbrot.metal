#include <metal_stdlib>
using namespace metal;

// Code to support mandelbrot set calculations.

constant int maxRainbow = 256 * 3;
constant int maxIterations = 250;

int section(int n) {
	return 256 * n;
}

struct RGBA {
  uint8_t r;
  uint8_t g;
  uint8_t b;
  uint8_t a;
};

struct RGBA colorFromIndex(int i) {
  i = i % section(5);
  // Red to yellow.
  if(i < section(1)) {
	  return RGBA{255, (uint8_t)i, 0, 255};
  }
  // Yellow to green.
  if(i < section(2)) {
	  return RGBA{(uint8_t)(section(2) - i - 1), 255, 0, 255};
  }
  // Green to light blue.
  if(i < section(3)) {
	  return RGBA{0, 255, (uint8_t)(section(2) + i), 255};
  }
  // Light blue to dark blue.
  if(i < section(4)) {
	  return RGBA{0, (uint8_t)(section(4) - i - 1), 255, 255};
  }
  // Dark blue to purple.
  return RGBA{(uint8_t)(section(4) + i), 0, 255, 255};
}

float scale(float fromMax, float toMin, float toMax, float v) {
	return ((v / fromMax) * (toMax - toMin)) + toMin;
}

// isInSet returns 0 for numbers that are in the set, or the number of iterations taken to escape.
int isInSet(float creal, float cimag) {
	float zreal = creal;
        float zimag = cimag;
	for(int n = 0; n < maxIterations; n++) {
                float zzreal = zreal;
                float zzimag = zimag;
		zreal = (zzreal*zzreal - zzimag*zzimag) + creal;
		zimag = (zzreal*zzimag + zzimag*zzreal) + cimag;
		if(zreal > 2.0 || zimag > 2.0 ) {
			return n;
		}
	}
	return 0;

}

// Normal metal code from here onwards.

typedef struct Params {
  int w_in, h_in, d_in;
  int w_out, h_out, d_out;
} Params;

int idx(int x, int y, int z, int w, int h, int d) {
  int i = z * w * h;
  i += y * w;
  i += x;
  return i;
}

kernel void process(device const Params* p,
    device uint8_t* input, 
    device uint8_t* output, 
    uint3 gridSize[[threads_per_grid]],
    uint3 gid[[thread_position_in_grid]]) {

  // Only process once per pixel of data (4 uint8_t)
  if(gid.x % 4 != 0) {
    return;
  }

  int x = gid.x / 4;
  int y = gid.y;
  int w = gridSize[0];
  int h = gridSize[1];
  
  // Since we know we're in the first column...
  // we can process the whole row.
  int index = idx(gid.x, gid.y, gid.z,
    p->w_in, p->h_in, p->d_in);

  // Parameters to define the visible area.
  float min_r = -1.4;
  float max_r = 3.0;
  float min_i = -0.8;
  float max_i = 0.8;

  // Show some numbers.
  float r = scale(w, min_r, max_r, float(x));
  float i = scale(h, min_i, max_i, float(y));
  int n = isInSet(r, i);
  if(n == 0) {
    output[index+0] = 0,
    output[index+1] = 0,
    output[index+2] = 0,
    output[index+3] = 255;
  } else {
    float rainbowIndex = scale(float(maxIterations), 0.0, float(maxRainbow), float(n));
    RGBA c = colorFromIndex(int(rainbowIndex));
    output[index+0] = c.r;
    output[index+1] = c.g;
    output[index+2] = c.b;
    output[index+3] = 255;
  }
}

