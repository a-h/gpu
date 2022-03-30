// +build darwin

#include "mtl.h"
#import <Metal/Metal.h>

id<MTLDevice> device;
id<MTLComputePipelineState> pipelineState;
id<MTLCommandQueue> commandQueue;

void setup(char *source) {
  device = MTLCreateSystemDefaultDevice();
  // NSLog(@"Using default device %s", [device.name UTF8String]);

  // Create library of code.
  NSError *error = nil;
  MTLCompileOptions *compileOptions = [MTLCompileOptions new];
  compileOptions.languageVersion = MTLLanguageVersion1_1;
  NSString *ss = [NSString stringWithUTF8String:source];
  id<MTLLibrary> newLibrary = [device newLibraryWithSource:ss
                                                   options:compileOptions
                                                     error:&error];
  if (newLibrary == nil) {
    NSLog(@"Failed to create new library, error %@.", error);
    return;
  }

  // Add the process function.
  id<MTLFunction> processFunction =
      [newLibrary newFunctionWithName:@"process"];
  if (processFunction == nil) {
    NSLog(@"Failed to find the process function.");
    return;
  }

  //NSLog(@"%@", [newLibrary functionNames]);

  // Create a compute pipeline state object.
  pipelineState = [device newComputePipelineStateWithFunction:processFunction
                                                        error:&error];
  if (pipelineState == nil) {
    NSLog(@"Failed to created pipeline state object, error %@.", error);
    return;
  }

  commandQueue = [device newCommandQueue];
  if (commandQueue == nil) {
    NSLog(@"Failed to find the command queue.");
    return;
  }
}

id<MTLBuffer> bufferInput;
id<MTLBuffer> bufferOutput;

void createBuffers(void* in, int in_data_size_bytes, int in_array_size, 
    void* out, int out_data_size_bytes, int out_array_size) {
  bufferInput = [device newBufferWithBytes:in 
                             length:in_array_size*in_data_size_bytes 
                            options:MTLResourceStorageModeShared];
  bufferOutput = [device newBufferWithBytes:out 
                             length:out_array_size*out_data_size_bytes 
                            options:MTLResourceStorageModeShared];
}

void *run(Params *params) {
  @autoreleasepool {
    NSError *error = nil;

    // Send compute command.
    id<MTLCommandBuffer> commandBuffer = [commandQueue commandBuffer];
    if (commandBuffer == nil) {
      NSLog(@"Failed to get the command buffer.");
      return nil;
    }
    // Get the compute encoder.
    id<MTLComputeCommandEncoder> computeEncoder =
        [commandBuffer computeCommandEncoder];
    if (computeEncoder == nil) {
      NSLog(@"Failed to get the compute encoder.");
      return nil;
    }

    // Create the data to pass in to the add function.
    // Buffers to hold data.
    // Encode the pipeline state object and its parameters.
    [computeEncoder setComputePipelineState:pipelineState];
    // The inputs.
    [computeEncoder setBytes:params length:24 atIndex:0]; // 24 bytes (32 bits * 6)
    [computeEncoder setBuffer:bufferInput offset:0 atIndex:1];
    [computeEncoder setBuffer:bufferOutput offset:0 atIndex:2];

    MTLSize threadsPerGrid = MTLSizeMake(params->w_in, params->h_in, params->d_in);

    // Calculate a threadgroup size.
    // https://developer.apple.com/documentation/metal/calculating_threadgroup_and_grid_sizes?language=objc
    NSUInteger w = pipelineState.threadExecutionWidth;
    NSUInteger h = pipelineState.maxTotalThreadsPerThreadgroup / w;
    MTLSize threadsPerThreadgroup = MTLSizeMake(w, h, 1);

    // Encode the compute command.
    [computeEncoder dispatchThreads:threadsPerGrid
              threadsPerThreadgroup:threadsPerThreadgroup];

    // End the compute pass.
    [computeEncoder endEncoding];

    // Execute the command.
    [commandBuffer commit];

    // Normally, you want to do other work in your app while the GPU is running,
    // but in this example, the code simply blocks until the calculation is
    // complete.
    [commandBuffer waitUntilCompleted];

    return bufferOutput.contents;
  }
}
