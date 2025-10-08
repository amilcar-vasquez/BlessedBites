//file: ui/static/js/app.js

document.addEventListener('DOMContentLoaded', function() {
    setTimeout(function() {
        var elems = document.querySelectorAll('select');
        M.FormSelect.init(elems);
    }, 100); // Give it 100ms to ensure DOM is stable

    // Initialize order from localStorage
    const existing = JSON.parse(localStorage.getItem('blessed_order') || '[]');
    syncOrderState(existing);
});


//debounces search
document.addEventListener("DOMContentLoaded", function () {
    const searchInput = document.getElementById("searchBox");
    const cardContainer = document.getElementById("menuCards");

    let debounceTimeout;

    searchInput.addEventListener("input", function () {
        const query = searchInput.value.trim();
        clearTimeout(debounceTimeout);

        if (query.length > 0) {
            debounceTimeout = setTimeout(() => {
                fetch(`/search.json?q=${encodeURIComponent(query)}`)
                    .then((res) => {
                        if (!res.ok) throw new Error("Network response was not ok");
                        return res.json();
                    })
                    .then((data) => {
                        if (cardContainer) {
                            if (data.length === 0) {
                                cardContainer.innerHTML = `<p class="center-align">No results found.</p>`;
                            } else {
                                cardContainer.innerHTML = data
                                    .map((item) => `
                                        <div class="col s12 m6 l4 menu-card">
                                            <div class="card">
                                                <div class="card-image">
                                                    <img class="responsive-img" src="${item.image_url}" alt="${item.name}">
                                                </div>
                                                <div class="card-content">
                                                    <span class="card-title">${item.name}</span>
                                                    <p>${item.description}</p>
                                                    <p><strong>$${item.price.toFixed(2)}</strong></p>
                                                </div>
                                                <div class="card-btn">
                                                    <a href="#">Add to Order</a>
                                                </div>
                                            </div>
                                        </div>
                                    `)
                                    .join('');
                            }
                        }
                    })
                    .catch((err) => {
                        console.error("Search error:", err);
                        if (cardContainer) {
                            cardContainer.innerHTML = "<p class='center-align'>Error loading search results.</p>";
                        }
                    });
            }, 300);
        } else {
            location.reload(); // Restore original card list if input is cleared
        }
    });
});


// Use event delegation to catch dynamically loaded menu items
document.addEventListener('click', function (e) {
    const target = e.target.closest('.btn-add-order');
    if (!target) return;
    e.preventDefault();

    const menuItemID = parseInt(target.dataset.id);
    const itemName = target.dataset.name;
    const itemPrice = parseFloat(target.dataset.price);

    // Prompt user for quantity
    const itemAmount = parseInt(prompt(`Enter quantity for ${itemName}:`, "1"), 10);
    if (isNaN(itemAmount) || itemAmount <= 0) {
        alert("Invalid quantity. Please enter a positive number.");
        return;
    }

    // Load existing order from localStorage
    let orderData = JSON.parse(localStorage.getItem('blessed_order') || '[]');
    orderData.push({ id: menuItemID, name: itemName, price: itemPrice, qty: itemAmount });
    localStorage.setItem('blessed_order', JSON.stringify(orderData));

    // Sync to any hidden inputs and render
    syncOrderState(orderData);
});

function renderOrder(orderData) {
    const orderList = document.getElementById('orderList');
    const orderTotal = document.getElementById('orderTotal');
    orderList.innerHTML = '';
    let total = 0;
    orderData.forEach(item => {
        const li = document.createElement('li');
        li.textContent = `${item.qty} x ${item.name} - $${(item.price * item.qty).toFixed(2)}`;
        orderList.appendChild(li);
        total += item.price * item.qty;
    });
    orderTotal.textContent = `Total: $${total.toFixed(2)}`;
    // Update cart badges (desktop and mobile)
    const count = orderData.reduce((acc, it) => acc + (it.qty || 0), 0);
    const desktopEl = document.getElementById('cartCount');
    if (desktopEl) {
        desktopEl.textContent = count;
        desktopEl.style.display = count > 0 ? 'inline-block' : 'none';
    }
    document.querySelectorAll('.cart-count').forEach(el => {
        el.textContent = count;
        el.style.display = count > 0 ? 'inline-block' : 'none';
    });
}

function syncOrderState(orderData) {
    // Render order UI if present
    try {
        renderOrder(orderData);
    } catch (err) {
        // ignore missing DOM elements
    }
    // Sync hidden inputs used by forms
    const hid = document.getElementById('orderData');
    if (hid) hid.value = JSON.stringify(orderData || []);
    const hid2 = document.getElementById('orderDataCheckout');
    if (hid2) hid2.value = JSON.stringify(orderData || []);
    // Update cart badges even if not rendering order list
    const count = (orderData || []).reduce((acc, it) => acc + (it.qty || 0), 0);
    const desktopEl2 = document.getElementById('cartCount');
    if (desktopEl2) {
        desktopEl2.textContent = count;
        desktopEl2.style.display = count > 0 ? 'inline-block' : 'none';
    }
    document.querySelectorAll('.cart-count').forEach(el => {
        el.textContent = count;
        el.style.display = count > 0 ? 'inline-block' : 'none';
    });
}

function flyToCart(elem) {
    const clone = elem.cloneNode(true);
    const rect = elem.getBoundingClientRect();
    clone.style.position = 'fixed';
    clone.style.left = `${rect.left}px`;
    clone.style.top = `${rect.top}px`;
    clone.classList.add('fly-anim');
    document.body.appendChild(clone);

    setTimeout(() => {
        clone.remove();
    }, 700);
}
// Example usage:
// flyToCart(document.querySelector('.dish-image'));
