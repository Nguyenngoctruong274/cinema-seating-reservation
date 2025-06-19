# 🎬 Cinema Seating Reservation System – Take-Home Assignment

## 🔒 Confidentiality Notice
This assignment is **confidential** and provided as part of your hiring process with **VulcanLabs**. Do **not share** or publish this assignment, including its content, code, design, or related discussions with anyone outside of this process.

---

## 🎯 Problem Statement: Cinema Seating with Social Distancing

To comply with public health regulations during movie screenings, cinema seating must follow **social distancing rules**:

- ✅ People in the **same group** can sit together without any distancing requirement.
- ❌ **Different groups** must maintain a **minimum Manhattan distance** from each other.
- 💡 Note: A group can consist of `one` or `multiple people`. People within the same group can sit next to each other without any distancing restriction. However, different groups must maintain the required minimum distance from each other.
---

## 📐 Manhattan Distance Formula

Given two seats at coordinates `(x1, y1)` and `(x2, y2)`, the **Manhattan distance** between them is:

```
Distance = |x1 - x2| + |y1 - y2|
```

---

## 🧮 Example Cinema Seating Layout

Consider a layout of `4 rows × 5 columns` with `min_distance = 4`:

```txt
Legend:
1 - reserved seat
0 - empty seat

Layout:
[ [1, 1, 0, 0, 0],
  [1, 1, 1, 0, 0],
  [1, 0, 0, 0, 0],
  [0, 0, 0, 0, 1] ]
```

The **shortest Manhattan distance** is between seat `(1, 2)` (part of the left group) and seat `(3, 4)` (a separate group):

```
|1 - 3| + |2 - 4| = 2 + 2 = 4
```

✅ With `min_distance = 4`, this setup is **valid**.

---

## ✅ Your Task

Create a **RESTful** or **gRPC** service to manage seating in the cinema while enforcing social distancing rules.

### 1. Configure Cinema Layout
- Inputs:
    - `rows` (int): total number of rows
    - `columns` (int): total number of columns
    - `min_distance` (int): minimum allowed Manhattan distance between different groups

### 2. Query Available Seats
- Input: `number_of_seats` (int)
- Output: a list of contiguous seat blocks that can be booked together as a group

### 3. Check Available Seats
- Input: list of `(row, column)` coordinates
- Output: a list of available seats

### 4. Reserve Seats
- Input: list of `(row, column)` coordinates
- Behavior: Reserve seats if the reservation does not violate the distance rule

### 5. Cancel Reservation
- Input: list of `(row, column)` coordinates
- Behavior: Frees up the seats for future reservations

---

## ⚙️ Additional Technical Requirements

### 🧵 Concurrency & Race Conditions
- Handle **simultaneous seat reservations** correctly
- Use atomic operations or locking to avoid double bookings

### 🛑 Deadlock Prevention
- Use consistent locking strategies or timeout-based retry
- Avoid blocking I/O in shared critical sections

### 🚦 High Traffic Scalability
- Assume thousands of users may attempt booking simultaneously
- Suggestions:
    - In-memory stores for fast access
    - Sharding by cinema or screen
    - Load-balanced APIs

### 🧪 Testability
- Cover edge cases:
    - All seats reserved
    - Overlapping requests
    - Invalid coordinates
    - Too many seats requested

---

## 💡 Business Assumptions

- Groups must reserve all their seats together (no partial allocations)
- Seat IDs are zero-indexed with `(0, 0)` at the **top-left**
- Each reservation creates a **new group**
- System may support multiple screens or theaters in the future

---

## 📦 Deliverables

- ✅ Source code
- 📄 `README.md`:
    - Setup instructions
    - API documentation or `.proto` file if gRPC
    - Assumptions and design decisions
- 🧪 Test cases (unit + optional integration)
- ❗ Keep this assignment confidential
