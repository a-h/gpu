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
    device const float* input, 
    device float* output, 
    uint3 gridSize[[threads_per_grid]],
    uint3 gid[[thread_position_in_grid]]) {
  // Only process once per row of data.
  if(gid.x != 0) {
    return;
  }

  // Since we know we're in the first column...
  // we can process the whole row.
  int input_index = idx(gid.x, gid.y, gid.z,
    p->w_in, p->h_in, p->d_in);

  float a = input[input_index];
  float b = input[input_index+1];

  int output_index = idx(0, gid.y, 0,
    p->w_out, p->h_out, p->d_out);

  output[output_index] = a + b;
}
