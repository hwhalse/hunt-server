const map = new maplibregl.Map({
  container: 'map',
  style: 'https://demotiles.maplibre.org/style.json',
  center: [-98.5795, 39.8283], // USA center
  zoom: 4
});

map.addControl(new maplibregl.NavigationControl());

const wsUrlEl = document.getElementById('wsUrl');
const connectBtn = document.getElementById('connect');
const disconnectBtn = document.getElementById('disconnect');
const deleteBtn = document.getElementById('deleteBtn');
const logEl = document.getElementById('log');
const statusDot = document.getElementById('statusDot');
const statusText = document.getElementById('statusText');

let ws = null;
let uid = "";
const events = {
    Init: 1,
    InitResponse: 2,
    CommandNodes: 3,
    Locations: 4,
    Units: 5,
    LocationUpdate: 6,
    CommandNodeUpdate: 7,
    TargetUpdate: 8,
    NewCommandNode: 9,
    NewGroup: 10,
    UnitUpdate: 11,
    Disconnect: 12,
    UpdateCallsign: 13,
    NewUid: 14,
    CommandNodeDelete: 15,
    CommandNodeStatusUpdate: 16,
    Error: 99
}

function log(...args) {
  const line = document.createElement('div');
  line.textContent = '[' + new Date().toLocaleTimeString() + '] ' + args.join(' ');
  logEl.appendChild(line);
  logEl.scrollTop = logEl.scrollHeight;
}

function setStatus(state) {
  const color = state === 'connected' ? '#0a0' : state === 'connecting' ? '#fa0' : state === 'error' ? '#a00' : '#888';
  statusText.textContent = state;
  statusDot.style.background = color;
}

function connect() {
  if (ws) ws.close();
  const url = wsUrlEl.value.trim();
  if (!url) return alert('Enter ws://localhost:PORT');

  setStatus('connecting');
  ws = new WebSocket(url);

  ws.addEventListener('open', () => {
    setStatus('connected');
    connectBtn.disabled = true;
    disconnectBtn.disabled = false;
    deleteBtn.disabled = false;
    sendInit()
    log('Connected to', url);
  });

  ws.addEventListener('message', ev => {
    log('IN:', ev.data);
    try {
      const data = JSON.parse(ev.data);
      switch (data.type) {
        case events.NewUid: {
            let received = data.payload;
            if (typeof received === "string") {
                try {
                    const maybe = JSON.parse(received);
                    if (typeof maybe === "string") received = maybe;
                } catch (e) {
                    log("Unable to parse incoming UID")
                }
            }
            log('NEW ID (parsed):', received);
            uid = received;
            break;
        }
        case events.InitResponse: {
            const payload = typeof data.payload === "string" ? JSON.parse(data.payload) : data.payload
            const locations = payload.users

            log("Incoming state", locations)

            if (!Array.isArray(locations)) {
                log("Invalid payload format:", payload)
                return
            }

            for (const userLocation of locations) {
                log("Location in locations", userLocation.location, userLocation.callsign)
                try {
                    const { location, callsign } = userLocation
                    if (!location || typeof location.lat !== "number" || typeof location.lon !== "number") {
                        log("Skipping invalid location:", userLocation)
                        continue
                    }

                    log('LOCATION+CALLSIGN:', location.lat, callsign)
                    if (location.lat == 0.0 && location.lon == 0.0) return
                    if (location.lat > -90 && location.lat < 90) {
                        addMarker(location.lat, location.lon, callsign)
                    }
                } catch (e) {
                    console.error("Error processing location:", e)
                }
            }
            break
}

        case events.Locations: {
            const locations = JSON.parse(data.payload)
            for (const userLocation of locations) {
                try {
                    const { location, callsign } = userLocation
                    log('LOCATION+CALLSIGN:', location.lat, callsign)
                
                    if (location.lat > -90 && location.lat < 90) {
                        addMarker(location.lat, location.lon, callsign)
                    }
                } catch(e) {
                    alert(e)
                }
            }
        }
        case events.LocationUpdate: {
            const update = JSON.parse(data.payload)
            updateMarkerPosition(update.callsign, update.location.lat, update.location.lon)
        }
      }
      log('EVENT:', data.type)
    } catch (e) {
        alert(e)
    }
  });

  ws.addEventListener('close', () => {
    setStatus('disconnected');
    connectBtn.disabled = false;
    disconnectBtn.disabled = true;
    deleteBtn.disabled = true;
    log('Disconnected');
  });

  ws.addEventListener('error', () => {
    setStatus('error');
    log('WebSocket error');
  });
}

function disconnect() {
  if (ws) ws.close();
}

function sendInit() {
  if (!ws || ws.readyState !== WebSocket.OPEN) return alert('Not connected');
  const init = {
    Type: 1,
    Payload: JSON.stringify({
        callsign: "tester"
    })
  }
  ws.send(JSON.stringify(init));
  log('OUT:', JSON.stringify(init));
}

const markers = {}

function updateMarkerPosition(name, lat, lng) {
    if (markers[name]) markers[name].setLngLat([lng, lat]).addTo(map)
    else addMarker(lat, lng, name)
}

function addMarker(lat, lng, name) {
    log('ADD MARKER', lat, lng, name)
    const el = document.createElement('div');
    el.style.background = 'red';
    el.style.width = '12px';
    el.style.height = '12px';
    el.style.borderRadius = '50%';
    el.style.border = '2px solid white';

    const marker = new maplibregl.Marker(el)
        .setLngLat([lng, lat])
        .addTo(map);
    markers[name] = marker
}

async function deleteLocations() {
    await fetch('http://localhost:8080/delete')
}

connectBtn.addEventListener('click', connect);
disconnectBtn.addEventListener('click', disconnect);
deleteBtn.addEventListener('click', deleteLocations)