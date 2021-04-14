def matrix_mul_ijk(x, y):
    z = [[0 for row in range(2)] for col in range(2)]

    for i in range(2):
        print("i", i)
        for j in range(2):
            print("j", j)
            total = 0

            for k in range(2):
                print(x[i][k], y[k][j])
                total += x[i][k] * y[k][j]

            z[i][j] = total

    print(z)
    return z

def matrix_mul_ikj(x, y):
    z = [[0 for row in range(2)] for col in range(2)]

    for i in range(2):
        print("i", i)
        for k in range(2):
            print("k", k)
            for j in range(2):
                print(x[i][k], y[k][j])
                z[i][j] += x[i][k] * y[k][j]

    print(z)
    return z



# assert matrix_mul_ijk([[7,8], [2,9]], [[14,5], [5, 18]]) == [[138, 179], [73, 172]]
assert matrix_mul_ikj([[7,8], [2,9]], [[14,5], [5, 18]]) == [[138, 179], [73, 172]]
