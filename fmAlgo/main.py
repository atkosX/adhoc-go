import hashlib
import matplotlib.pyplot as plt
import numpy as np

def rightmost_bit_pos(x):
    if x == 0:
        return -1
    pos = 0
    while (x & 1) == 0:
        pos += 1
        x >>= 1
    return pos

def hash_item(item, seed):
    s = f"{item}-{seed}".encode()
    h = hashlib.sha256(s).digest()
    return int.from_bytes(h[:8], byteorder='big')

def fm_algo(items, num_hashes=32):
    estimates = []
    for i in range(num_hashes):
        max_rho = 0
        for item in items:
            h = hash_item(item, i)
            rho = rightmost_bit_pos(h)
            max_rho = max(max_rho, rho)
        estimates.append(1 << max_rho)
    estimates.sort()
    median = estimates[num_hashes // 2]
    correction = 0.77351
    return int(median / correction)

# Test
actual = []
estimated = []
for n in range(100, 2000, 100):
    items = [f"item_{i}" for i in range(n)]
    est = fm_algo(items, num_hashes=256)
    actual.append(n)
    estimated.append(est)

# Plot
plt.figure(figsize=(10, 6))
plt.plot(actual, actual, label="Actual", linestyle="--", color="blue")
plt.plot(actual, estimated, 'o-', label="FM Estimate", color="red", markersize=4)
plt.xlabel("Actual Unique Items")
plt.ylabel("Estimated Unique Items")
plt.title("Flajolet-Martin Cardinality Estimation")
plt.legend()
plt.grid(True)
plt.tight_layout()

# Save the plot instead of showing it
plt.savefig("fm_estimation_plot.png", dpi=300, bbox_inches='tight')
print("Plot saved as 'fm_estimation_plot.png'")

# Print some statistics
errors = [abs(a - e) / a * 100 for a, e in zip(actual, estimated)]
avg_error = sum(errors) / len(errors)
print(f"Average estimation error: {avg_error:.2f}%")
print(f"Min error: {min(errors):.2f}%")
print(f"Max error: {max(errors):.2f}%")

# Show a few examples
print("\nSample results:")
for i in range(0, len(actual), 3):
    print(f"Actual: {actual[i]:4d}, Estimated: {estimated[i]:4d}, Error: {errors[i]:5.1f}%")
