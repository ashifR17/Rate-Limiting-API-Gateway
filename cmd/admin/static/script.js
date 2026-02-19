const API_BASE = 'http://localhost:9090/admin';

const loadBtn = document.getElementById('loadBtn');
const updateBtn = document.getElementById('updateBtn');

// Dummy API endpoints
const apiEndpoints = [
    "/api/health",
    "/api/ping",
    "/api/users",
    "/api/orders",
    "/api/products",
];

const endpointList = document.getElementById("endpointList");

// Render endpoint checkboxes
function renderEndpoints(savedConfig = {}) {
    endpointList.innerHTML = "";
    apiEndpoints.forEach(ep => {
        const div = document.createElement("div");
        div.className = "endpoint-item";

        const input = document.createElement("input");
        input.type = "checkbox";
        input.id = `endpoint-${ep.replace(/\//g, "-")}`;
        input.checked = savedConfig[ep] ?? true; // default enabled

        const label = document.createElement("label");
        label.htmlFor = input.id;
        label.innerText = ep;

        div.appendChild(input);
        div.appendChild(label);
        endpointList.appendChild(div);
    });
}

// Initialize dummy endpoints (all enabled by default)
renderEndpoints();

// Load and Update Config buttons
loadBtn.addEventListener('click', () => fetchConfig(true));
updateBtn.addEventListener('click', updateConfig);

// Save Endpoint config
document.getElementById("saveEndpointsBtn").addEventListener("click", () => {
    const token = document.getElementById('adminKey').value.trim();
    if (!token) return alert("Enter admin token");

    const configToSave = {};
    apiEndpoints.forEach(ep => {
        const input = document.getElementById(`endpoint-${ep.replace(/\//g, "-")}`);
        configToSave[ep] = input.checked;
    });

    console.log("Dummy endpoint config to save:", configToSave);
    alert("Endpoint configuration saved (dummy, backend not connected yet)");
    // Later: send to backend using fetch()
});

async function fetchConfig(showAlert = false) {
    const token = document.getElementById('adminKey').value.trim();
    if (!token) {
        alert("Please enter admin token.");
        return;
    }

    try {
        const res = await fetch(`${API_BASE}/config`, {
            headers: { 'X-Admin-Key': token }
        });

        if (!res.ok) throw new Error("Failed to fetch config");

        const cfg = await res.json();

        document.getElementById('globalCapacity').value = cfg.global_capacity;
        document.getElementById('globalRate').value = cfg.global_rate;
        document.getElementById('userCapacity').value = cfg.user_capacity;
        document.getElementById('userRate').value = cfg.user_rate;
        document.getElementById('userApiCapacity').value = cfg.user_api_capacity;
        document.getElementById('userApiRate').value = cfg.user_api_rate;

        // TODO: load API endpoints if available from backend
    } catch (err) {
        console.error(err);
        if (showAlert) alert("Error fetching config. Check admin token or server.");
    }
}

async function updateConfig() {
    const token = document.getElementById('adminKey').value.trim();
    if (!token) return alert("Please enter admin token.");

    const cfg = {
        global_capacity: Number(document.getElementById('globalCapacity').value),
        global_rate: Number(document.getElementById('globalRate').value),
        user_capacity: Number(document.getElementById('userCapacity').value),
        user_rate: Number(document.getElementById('userRate').value),
        user_api_capacity: Number(document.getElementById('userApiCapacity').value),
        user_api_rate: Number(document.getElementById('userApiRate').value)
    };

    try {
        const res = await fetch(`${API_BASE}/config`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Admin-Key': token
            },
            body: JSON.stringify(cfg)
        });

        if (!res.ok) throw new Error("Failed to update config");

        alert("Rate limit config updated successfully.");
        fetchConfig(true);
    } catch (err) {
        console.error(err);
        alert("Error updating config.");
    }
}
