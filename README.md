# **TriFS: A Distributed File System Optimized for Tiny Files**

## **Introduction**

TriFS is a distributed file system (DFS) specifically engineered to address the challenges of storing and managing massive volumes of tiny files (typically \< 64KB, but effectively up to \~1MB). Traditional DFS designs, optimized for large, sequential data workloads, exhibit significant inefficiencies when faced with billions or trillions of small, discrete data records. TriFS is built from the ground up in Golang to provide a scalable, performant, and storage-efficient solution for these modern data patterns.

## **The Problem with Tiny Files**

Conventional DFS architectures struggle with tiny files due to:

* **Metadata Bottleneck:** The sheer volume of metadata required for billions of files overwhelms centralized metadata servers.  
* **Storage Inefficiency:** Large, fixed-size data blocks lead to massive internal fragmentation when storing small files.  
* **High I/O Overhead:** The fixed cost of network communication and server processing per operation becomes dominant for frequent accesses to small files.

TriFS is designed to mitigate these issues.

## **Key Features**

* **Optimized for Tiny Files:** Core design decisions prioritize efficient handling of small data records.  
* **Distributed Metadata:** Scales metadata management horizontally using a distributed Key-Value store.  
* **Data Packing:** Efficiently utilizes storage by packing multiple tiny files into larger chunks.  
* **Copy-on-Write (COW) Mutations:** Simplifies updates and appends while supporting versioning.  
* **Tombstone Deletes:** Provides fast logical deletion with background physical removal.  
* **Background Garbage Collection & Compaction:** Reclaims wasted space and improves read efficiency over time.  
* **Erasure Coding:** Ensures data fault tolerance and durability with better storage efficiency than simple replication.  
* **Thick Client:** Optimizes performance through aggressive caching and operation batching.  
* **Virtual Hierarchy:** Presents a familiar file system tree structure to the user via the client API, independent of physical data layout.  
* **Snapshotting:** Enables point-in-time views of the file system state.  
* **Compression:** Reduces storage footprint and I/O using fast algorithms.

## **Architecture Overview**

TriFS employs a distributed architecture with a clear separation of concerns:

* **Metadata System:** A distributed Key-Value store manages file metadata (File ID, location, size, version, checksum, etc.) and the virtual directory structure.  
* **Data Servers:** Store data in fixed-size chunks (\~10MB), packing multiple tiny files within each chunk. These chunks are erasure coded across servers for fault tolerance.  
* **Client:** The application-facing component that interacts with the Metadata and Data Servers. It manages caching, batches requests, and presents a virtual file system hierarchy to the user.

For a more detailed architectural breakdown, refer to the [TriFS Design Document](https://docs.google.com/document/d/161QHUgER5yCfzgVeeZBj3hUqnxjkTUcs90stOqHaAPo/edit?usp=sharing).

## **How it Works (High-Level)**

1. **Write:** Client requests a write. Metadata system allocates a location (Chunk ID). Client sends data to the appropriate Data Server(s) for packing and erasure coding. Metadata system records the file's location and attributes.  
2. **Read:** Client requests a file using its path. Client's virtual hierarchy logic translates path to File ID via metadata lookup (potentially cached). Client uses File ID to determine Chunk ID and Data Server(s). Client reads the relevant portion of the chunk from the Data Server(s).  
3. **Update/Append:** New data is written using COW to a new location. Metadata is updated to point to the new version.  
4. **Delete:** Metadata entry is marked with a tombstone.  
5. **Background:** GC/Compaction processes run periodically to clean up obsolete data from chunks and reclaim space.

## **Use Cases**

TriFS is ideal for workloads characterized by massive numbers of small files and high I/O rates on these files, including:

* Storing and serving IoT sensor data.  
* Managing application and system log files.  
* Hosting large collections of small web assets (images, CSS, JS).  
* Storing configuration files for distributed systems.  
* Archiving and querying segments of time-series data.
