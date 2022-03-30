#include <metal_stdlib>
using namespace metal;

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

  // Since we know we're in the first column...
  // we can process the whole row.
  int input_index = idx(gid.x, gid.y, gid.z,
    p->w_in, p->h_in, p->d_in);

  uint8_t r = input[input_index+0];
  uint8_t g = input[input_index+1];
  uint8_t b = input[input_index+2];
  uint8_t a = input[input_index+3];

  uint8_t avg = uint8_t((int(r) + int(g) + int(b)) / 3);

  output[input_index+0] = avg;
  output[input_index+1] = avg;
  output[input_index+2] = avg;
  output[input_index+3] = 255;
}
