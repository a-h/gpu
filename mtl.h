// +build darwin

typedef unsigned long uint_t;
typedef unsigned char uint8_t;
typedef unsigned short uint16_t;
typedef unsigned long long uint64_t;

// Matches with Params type in mtl.go
typedef struct Params {
  int w_in, h_in, d_in;
  int w_out, h_out, d_out;
} Params;

void compile(char* source);
void createBuffers(void* in, int in_data_size_bytes, int in_array_size,
    int out_data_size_bytes, int out_array_size);
void* run(Params *params);
