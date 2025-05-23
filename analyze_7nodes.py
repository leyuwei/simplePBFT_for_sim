# %%
import re
import numpy as np
import matplotlib.pyplot as plt
import numpy as np
import matplotlib.pyplot as plt
import matplotlib as mpl
from scipy.signal import find_peaks, peak_widths
from scipy.interpolate import interp1d
from scipy.signal import find_peaks, peak_prominences
from scipy.interpolate import UnivariateSpline
from scipy.optimize import curve_fit, minimize
from scipy.special import comb, factorial

plt.close('all')


mpl.rcParams['text.usetex'] = True
mpl.rcParams['font.family'] = 'serif'
mpl.rcParams['ps.usedistiller'] = 'ghostscript'

# Read the file
with open('6res.txt', 'r') as f:
    lines = f.readlines()

values = []

for line in lines:
    match = re.search(r'StoCLT:\{(\d+)\}', line)
    if match:
        num = float(match.group(1)) / 1000000
        if num < 800:
            continue
        values.append(num)
    # break

# Calculate the histogram
hist_values, bin_edges = np.histogram(values, bins=300, density=True)
bin_centers = (bin_edges[:-1] + bin_edges[1:]) / 2

# Find the significant peaks by adjusting the 'prominence' parameter
peaks, _ = find_peaks(hist_values, prominence=0.0005)  # Adjust prominence as needed
print(peaks)

# Define the PDF of the m-th order statistic from N i.i.d. exponential random variables
def order_statistic_pdf(y_real, theta, diff):
    N = 7  # Round N to nearest integer
    m = 5  # Round N to nearest integer
    y = y_real - diff
    # term1 = (theta**m * np.exp(-theta * y) * y**(m-1)) / factorial(m-1)
    # term2 = comb(N, m) * (1 - np.exp(-theta * y))**(N-m)

    term = comb(N-1, m) * m * (1 - np.exp(-theta * y))**(m-1) * np.exp(- theta * (N - m - 1) * y) * theta * np.exp(-theta * y)
    # term = np.exp(- (x + np.exp(-x))) # Gumbel Dist.
    return term

# Objective function to minimize: Sum of squared residuals
def objective(params):
    theta, diff = params
    print(params)
    predicted = order_statistic_pdf(bin_centers[peaks], theta, diff)
    residuals = hist_values[peaks] - predicted
    return np.sum(residuals**2)

# You need to provide an initial guess for the theta parameter
bp = 800
initial_guess = [0.1,bp]

# Fit the function to the data (using only the peaks for fitting)
result = minimize(objective, x0=initial_guess, bounds=[(0.00001, None), (bp, None)], method='Nelder-Mead', options={'maxiter': 10000, 'xatol': 1e-8, 'fatol': 1e-8})

# Generate dense x values for a smooth curve
dense_bin_centers = np.linspace(bin_centers.min(), bin_centers.max(), 1000)
dense_bin_centers = bin_centers

# Check the result
if result.success:
    fitted_params = result.x
    print(f"Fitted θ: {fitted_params[0]}")
    print(f"Fitted 1: {fitted_params[1]}")
else:
    print("Optimization failed:", result.message)

# Evaluate the fitted function with the optimized parameters
fitted_curve = order_statistic_pdf(dense_bin_centers, fitted_params[0], fitted_params[1])
for i in range(len(fitted_curve)):
    fitted_curve[i] = 0 if fitted_curve[i] < 0 else fitted_curve[i]

# Mean of the observed frequencies
mean_observed = np.mean(hist_values)
# Total sum of squares
sst = np.sum((hist_values - mean_observed) ** 2)
# Residual sum of squares
ssr = np.sum((hist_values - fitted_curve) ** 2)
# R^2 score
r_squared = 1 - ssr / sst
print('R^2 score:', r_squared)

plt.rcParams['font.family'] = 'Times New Roman'
plt.rcParams['font.size'] = 16

# Plotting
fig, ax = plt.subplots()
ax.set_axisbelow(True)
ax.grid(True, zorder=4)

ax.bar(bin_centers, hist_values, width=bin_edges[1] - bin_edges[0])
# plt.plot(bin_centers, hist_values, label='Histogram')
# plt.plot(bin_centers, spline(bin_centers), label='Smooth Envelope', linestyle='--')
# ax.scatter(bin_centers[peaks], hist_values[peaks], color='blue')  # Mark the significant peaks
# plt.plot(bin_centers[peaks], hist_values[peaks], "--", color='red')  # Mark the significant peaks
ax.plot(dense_bin_centers, fitted_curve, "-", color='red',linewidth=3.0)



## bandwidth to 30 megabits per second (Mbps) and ran the prototype on 4 fully connected wireless devices. size of a status report=320Kb
## The size of a status report was set to 40 kilobytes (kB), and the block size was set to 4 reports, corresponding to the 4 devices in the network
textbox_text = r'status report size=1.4 Mb' + '\n' + r'total number of CESs $M$=7' + '\n\n' + r'Fitted $\theta$=9.72 (blocks/second)' + '\n' + r'Fitted $\theta^{-1}_\mathrm{res}$=0.914 (seconds)' + '\n\n' + r'$R^2$ score=95.3\%'
plt.rc('text', usetex=True)
plt.rc('font', family='serif')
plt.text(0.95, 0.95, textbox_text, transform=plt.gca().transAxes,
         fontsize=16, verticalalignment='top', horizontalalignment='right',
         bbox=dict(facecolor='none', alpha=0.5, edgecolor='none'))

plt.xlim([900,2400])

plt.legend(["Analytical results", "Experimental results"],loc='lower right', fontsize=14)

plt.ylabel(r'PDF $f_{C_{q,z}}(t)$')
plt.xlabel('PBFT consensus delay (millisecond)')
plt.tight_layout()

plt.savefig('prototype-delaydist-2.eps', format='eps')

plt.show()