console.log('%c My Home project of Golang App', 'background: #0d6efd; color: #fff; padding:20px; margin:10px')
let folderIcon = '<svg width="100%" height="100%" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">\n' +
    ' <path d="M13 7L11.8845 4.76892C11.5634 4.1268 11.4029 3.80573 11.1634 3.57116C10.9516 3.36373 10.6963 3.20597 10.4161 3.10931C10.0992 3 9.74021 3 9.02229 3H5.2C4.0799 3 3.51984 3 3.09202 3.21799C2.71569 3.40973 2.40973 3.71569 2.21799 4.09202C2 4.51984 2 5.0799 2 6.2V7M2 7H17.2C18.8802 7 19.7202 7 20.362 7.32698C20.9265 7.6146 21.3854 8.07354 21.673 8.63803C22 9.27976 22 10.1198 22 11.8V16.2C22 17.8802 22 18.7202 21.673 19.362C21.3854 19.9265 20.9265 20.3854 20.362 20.673C19.7202 21 18.8802 21 17.2 21H6.8C5.11984 21 4.27976 21 3.63803 20.673C3.07354 20.3854 2.6146 19.9265 2.32698 19.362C2 18.7202 2 17.8802 2 16.2V7Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>\n' +
    ' </svg>'

let fileIcon = '<svg width="100%" height="100%" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">\n' +
    ' <path d="M14 2.26946V6.4C14 6.96005 14 7.24008 14.109 7.45399C14.2049 7.64215 14.3578 7.79513 14.546 7.89101C14.7599 8 15.0399 8 15.6 8H19.7305M20 9.98822V17.2C20 18.8802 20 19.7202 19.673 20.362C19.3854 20.9265 18.9265 21.3854 18.362 21.673C17.7202 22 16.8802 22 15.2 22H8.8C7.11984 22 6.27976 22 5.63803 21.673C5.07354 21.3854 4.6146 20.9265 4.32698 20.362C4 19.7202 4 18.8802 4 17.2V6.8C4 5.11984 4 4.27976 4.32698 3.63803C4.6146 3.07354 5.07354 2.6146 5.63803 2.32698C6.27976 2 7.11984 2 8.8 2H12.0118C12.7455 2 13.1124 2 13.4577 2.08289C13.7638 2.15638 14.0564 2.27759 14.3249 2.44208C14.6276 2.6276 14.887 2.88703 15.4059 3.40589L18.5941 6.59411C19.113 7.11297 19.3724 7.3724 19.5579 7.67515C19.7224 7.94356 19.8436 8.2362 19.9171 8.5423C20 8.88757 20 9.25445 20 9.98822Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>\n' +
    ' </svg>'

let currentIcon = fileIcon

let name = ''

async function loadAndRenderData(url, containerId) {
    try {
        // Получаем данные с указанного URL
        const response = await fetch(url);

        // Проверяем статус ответа (200–299 — успех)
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        // Парсим JSON
        const data = await response.json();

        // Находим контейнер для вывода
        const container = document.getElementById("FileManagerList");
        if (!container) {
            throw new Error('Container element not found');
        }

        // Генерируем HTML на основе данных
        let html = '';
        data.forEach(item => {
            if (item.is_dir == true){
                currentIcon = folderIcon
                name = `<td><a href="">${item.name || 'No name'}</a></td>`
            }
            else {
                currentIcon = fileIcon
                name = `<td>${item.name || 'No name'}</td>`
            }
            html += `<tr><td><div class="icon-file">${currentIcon}</div></td> 
${name}
 <td>${item.size_readable || '0'}</td>
<td>${item.created_at || '0'}</td> <td>text</td></tr>

     `;
        });
        html += '';

        // Вставляем HTML в контейнер
        container.innerHTML = html;
    } catch (error) {
        console.error('Error loading data:', error);
        document.getElementById(containerId).innerHTML =
            `<div class="error">Ошибка загрузки данных: ${error.message}</div>`;
    }
}

// Использование
loadAndRenderData('http://localhost:4200/api/disk?path=/DEV/my_home_server/', 'data-container');