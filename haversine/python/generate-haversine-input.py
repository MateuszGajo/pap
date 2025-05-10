import json
import random

def generate_pairs(num_records, filename):
    def random_coord():
        lon = random.uniform(-180, 180)
        lat = random.uniform(-90, 90)
        return lon, lat

    with open(filename, 'w') as f:
        f.write('{"pairs":[\n')
        for i in range(num_records):
            x0, y0 = random_coord()
            x1 = x0 + random.uniform(-0.1, 0.1)
            y1 = y0 + random.uniform(-0.1, 0.1)
            obj = json.dumps({"x0": round(x0, 6), "y0": round(y0, 6),
                              "x1": round(x1, 6), "y1": round(y1, 6)})
            if i < num_records - 1:
                f.write(obj + ",\n")
            else:
                f.write(obj + "\n")
        f.write(']}')

generate_pairs(10_000_000, 'data_10000000_flex.json')
