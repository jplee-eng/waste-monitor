function updateUI(reading) {
  document.getElementById("levelValue").textContent = reading.level;
  document.getElementById("batteryValue").textContent =
    reading.battery.toFixed(2);
  const rssi = reading.rssi;
  let signalStatus, signalClass;
  if (rssi >= -85) {
    signalStatus = "Excellent";
    signalClass = "good";
  } else if (rssi >= -95) {
    signalStatus = "Good";
    signalClass = "good";
  } else if (rssi >= -105) {
    signalStatus = "Fair";
    signalClass = "warning";
  } else {
    signalStatus = "Poor";
    signalClass = "error";
  }
  document.getElementById(
    "signalValue"
  ).innerHTML = `<span class="status ${signalClass}">${signalStatus}</span>`;
  let batteryClass =
    reading.battery > 3.3
      ? "good"
      : reading.battery > 3.0
      ? "warning"
      : "error";
  document.getElementById(
    "batteryValue"
  ).className = `metric status ${batteryClass}`;
  document.getElementById(
    "lastUpdate"
  ).textContent = `Last updated: ${new Date().toLocaleTimeString()}`;
  const statusHtml = `<p>RSSI: ${
    reading.rssi
  } dBm</p><p>Last Reading: ${new Date(
    reading.timestamp
  ).toLocaleString()}</p>`;
  document.getElementById("status").innerHTML = statusHtml;
}
const evtSource = new EventSource("/events");
evtSource.onmessage = (event) => {
  const reading = JSON.parse(event.data);
  updateUI(reading);
};
evtSource.onerror = () => {
  document.getElementById("status").innerHTML =
    '<p class="status error">Connection lost. Reconnecting...</p>';
};
