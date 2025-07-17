# **TriFS: A Distributed File System Optimized for Tiny Files**

> [!WARNING]
> TriFS is currently under active development. The foundational components are being built, and while the architecture is designed for scalability and robustness, this implementation is **not yet production-ready**. Explorers are highly encouraged to use the system and share their valuable feedback to help shape its future.

## **Introduction**

TriFS is a distributed file system (DFS) specifically engineered to address the challenges of storing and managing massive volumes of tiny files (typically \< 64KB, but effectively up to \~1MB). Traditional DFS designs, optimized for large, sequential data workloads, exhibit significant inefficiencies when faced with billions or trillions of small, discrete data records. TriFS is built from the ground up in Golang to provide a scalable, performant, and storage-efficient solution for these modern data patterns.

## **The Problem with Tiny Files**

Conventional DFS architectures struggle with tiny files due to:

- **Metadata Bottleneck:** The sheer volume of metadata required for billions of files overwhelms centralized metadata servers.
- **Storage Inefficiency:** Large, fixed-size data blocks lead to massive internal fragmentation when storing small files.
- **High I/O Overhead:** The fixed cost of network communication and server processing per operation becomes dominant for frequent accesses to small files.

TriFS is designed to mitigate these issues.

## **Key Features**

- **Optimized for Tiny Files:** Core design decisions prioritize efficient handling of small data records.
- **Distributed Metadata:** Scales metadata management horizontally using a **Distributed Metadata Service**.
- **Data Packing:** Efficiently utilizes storage by **packing multiple tiny files into larger, fixed-size data packs** (e.g., ~50MB). This tries to avoids internal fragmentation.
- **Copy-on-Write (COW) Mutations:** Simplifies updates and appends while supporting versioning.
- **Tombstone Deletes:** Provides fast logical deletion with background physical removal.
- **Background Garbage Collection & Compaction:** Reclaims wasted space and improves read efficiency over time by reorganizing data within packs.
- **Erasure Coding:** Ensures data fault tolerance and durability for **data packs** with better storage efficiency than simple replication.
- **Thick Client:** Optimizes performance through aggressive caching and operation batching.
- **Virtual Hierarchy:** Presents a familiar file system tree structure to the user via the client API, independent of physical data layout.
- **Snapshotting:** Enables point-in-time views of the file system state.
- **Compression:** Reduces storage footprint and I/O using fast algorithms.

## **Architecture Overview**

TriFS has a distributed architecture with a clear separation of concerns, primarily consisting of three main components:

- **Master:** Manages the overall file system metadata and coordinates operations. It does not store all the metadata directly, but rather points the Client to the appropriate Worker(s).
- **Worker:** Stores the actual file data in fixed-size **packs** (~50MB), efficiently packing multiple tiny files within each. These packs are erasure coded across Workers for fault tolerance.
- **Client:** The application-facing component that interacts with the **Master** for coordination and directly with **Workers** for data transfer. It manages caching, batches requests, and presents a virtual file system hierarchy to the user.

For a more detailed architectural breakdown, refer to the [TriFS Design Document](https://docs.google.com/document/d/161QHUgER5yCfzgVeeZBj3hUqnxjkTUcs90stOqHaAPo/edit?usp=sharing).

## **How it Works (High-Level)**

1.  **Write:** **Client** requests a write. **Master** allocates a location by choosing suitable **Worker(s)**. **Client** sends the data to the designated **Worker(s)** for packing and storage. **Master** records the file's basic information and attributes.
2.  **Read:** **Client** requests a file using its path. **Client** queries the **Master** to determine which **Worker(s)** hold the file's data. **Client** then reads the relevant portion of the pack directly from the appropriate **Worker(s)**.
3.  **Update/Append:** New data is written using Copy-on-Write principles to a new location (within a new or existing pack). Metadata managed by the **Master** is updated to point to the new version.
4.  **Delete:** The **Master** marks the file entry as deleted (a "tombstone").
5.  **Background:** Garbage collection and compaction processes run periodically on **Workers** to clean up obsolete data from **packs** and reclaim space.

## **Use Cases**

TriFS is ideal for workloads characterized by massive numbers of small files and high I/O rates on these files, including:

- Storing and serving IoT sensor data.
- Managing application and system log files.
- Hosting large collections of small web assets (images, CSS, JS).
- Storing configuration files for distributed systems.
- Archiving and querying segments of time-series data.
