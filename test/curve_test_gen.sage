import json
import os

# Convert points to affine coordinates if they are not at infinity
def point_to_coords(P):
    if P == E(0):
        return None 
    else:
        x, y = P.xy()
        return { "x": int(x), "y": int(y) }

def generate_point_not_on_curve(E, F):
    x = F.random_element()
    y_maybe_squared = x^3 + E.a4()*x + E.a6() # y^2 = x^3 + ax + b
    # if y_maybe_squared is indeed quadriatic residue then we do have a point on the curve
    # this means that we'll have to modify y to get a point which is definitely not on the curve
    if y_maybe_squared.is_square():
        y = y_maybe_squared.sqrt()
        # Modify y to ensure the point is not on the curve,
        # Let's say we'll modfiy y  as follows y' = y + 1 (mod p)
        # y'^2 = (y + 1)^2 = y^2 + 2y + 1
        # And since by definition y^2 = x^3 + ax + b, we have:
        # y'^2 = (x^3 + ax + b) + 2y + 1
        # So if 2y + 1 is 0 we are landing on the curve, thus we need to modify y differently
        # Note if we expand further we get y = - 1 / 2 but we are modding with p, so we can just check if y = p / 2  
        if y == F.order() // 2:
            y = (y + 2) % F.order()
        else:
            y = (y + 1) % F.order()

    # if y_maybe_squared is not a quadratic residue then we do not have a point on the curve,
    # so we can return any random value for y (since for given x we cannot be on the curve)
    else:
        y = F.random_element()
    return {"x": int(x), "y": int(y)}

N = 100 # Number of points to generate
p = random_prime(10000, lbound=1)
F = GF(p)

# Randomly select coefficients a and b ensuring the discriminant is non-zero
a = F.random_element()
b = F.random_element()
while 4*a^3 + 27*b^2 == 0:
    a = F.random_element()
    b = F.random_element()

E = EllipticCurve(F, [a,b])


# Generate N data points on the curve and store them in a list
points_on_curve = []
for _ in range(N):
    P = E.random_point()
    P_affine = point_to_coords(P)
    if P_affine is not None:
        points_on_curve.append(P_affine)

points_not_on_curve = []
for _ in range(N):
    point = generate_point_not_on_curve(E, F)
    points_not_on_curve.append(point)

output_data = {
    "p": int(p),
    "a": int(a),
    "b": int(b),
    "points_on_curve": points_on_curve,
    "points_not_on_curve": points_not_on_curve
}

json_data_on_curve = json.dumps(output_data, indent=2)

file_path = os.path.expanduser('~/source/weierstrass/test/curve_test_gen.json')
with open(file_path, 'w') as file:
    file.write(json_data_on_curve)
