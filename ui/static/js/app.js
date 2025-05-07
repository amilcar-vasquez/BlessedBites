//file: ui/static/js/app.js

document.addEventListener('DOMContentLoaded', function() {
    setTimeout(function() {
        var elems = document.querySelectorAll('select');
        M.FormSelect.init(elems);
    }, 100); // Give it 100ms to ensure DOM is stable
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
