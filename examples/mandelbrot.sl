import "../lib/std"

fun main() {
  mandel(-2.3, -1.3, 0.05, 0.07)
}

fun mandel(realstart, imagstart, realmag, imagmag: float) {
  mandelHelp(realstart, realstart+realmag*78.0, realmag,
             imagstart, imagstart+imagmag*48.0, imagmag)
}

fun mandelHelp(xMin, xMax, xStep, yMin, yMax, yStep: float) {
    for var y = yMin; y < yMax; y += yStep {
    for var x = xMin; x < xMax; x += xStep {
      var d = mandelConverge(x, y)
      printDensity(d)
    }
    print("\n")
  }
}

fun mandelConverge(real, imag: float): float {
  return mandelConverger(real, imag, 0.0, real, imag)
}

fun mandelConverger(real, imag, iters, creal, cimag: float): float {
  if iters > 255.0 {
    return iters;
  }
  if real*real + imag*imag > 4.0 {
    return iters;
  }
  
  return mandelConverger(
    real*real - imag*imag + creal,
    2.0*real*imag + cimag,
    iters + 1.0, creal, cimag)
}

fun printDensity(d: float) {
  if d > 8.0 {
    print(" ")
  } else if d > 4.0 {
    print(".")
  } else if d > 2.0 {
    print("+")
  } else {
    print("*")
  }
}
