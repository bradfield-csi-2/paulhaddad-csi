# Malloc and Free

* Hold memory allocation in chunks of multiples of 2 to improve cache
  utilization.
* Model the allocation as a contiguous block of memory
* Allocating a new block of memory
  * Find smallest continguous block of memory that holds the data
  * If that's not possible, then allocate it across the smallest number of
    contiguous chunks
  * Keep a separate data element that allows us to quickly determine if a chunk
    of memory is filled.
* Freeing memory
  * Indicate that the chunk is free
  * Check if we can consolidate 
