# gRPC-Observability
### Test JSON Data for RPCs

---

#### 1. **CreateUser**:
```json
{
  "user": {
    "name": "John Doe",
    "age": 25,
    "commuteMethod": "Car",
    "college": "Harvard University",
    "hobbies": "Reading"
  }
}
```

#### 2. **GetUser**:
```json
{
  "name": "John Doe"
}

```

#### 3. **UpdateUser**:
```json
{
  "name": "John Doe",
  "user": {
    "name": "John Doe",
    "age": 26,
    "commuteMethod": "Bike",
    "college": "Harvard University",
    "hobbies": "Reading, Cycling"
  }
}
```

#### 4. **DeleteUser**:
```json
{
  "name": "John Doe"
}

```
