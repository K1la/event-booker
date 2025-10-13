// Базовый URL API
const API_BASE = '/api/events';

// Утилиты для работы с API
class EventBookerAPI {
    constructor() {
        this.baseURL = API_BASE;
    }

    // Общий метод для выполнения HTTP запросов
    async request(url, options = {}) {
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || `HTTP error! status: ${response.status}`);
            }

            return data.result;
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    // Получить все мероприятия
    async getEvents() {
        return this.request(this.baseURL);
    }

    // Получить мероприятие по ID
    async getEventById(id) {
        return this.request(`${this.baseURL}/${id}`);
    }

    // Создать мероприятие
    async createEvent(eventData) {
        return this.request(this.baseURL, {
            method: 'POST',
            body: JSON.stringify(eventData)
        });
    }

    // Забронировать место
    async bookEvent(eventId, bookingData) {
        return this.request(`${this.baseURL}/${eventId}/book`, {
            method: 'POST',
            body: JSON.stringify(bookingData)
        });
    }

    // Подтвердить бронирование
    async confirmBooking(eventId) {
        return this.request(`${this.baseURL}/${eventId}/confirm`, {
            method: 'POST'
        });
    }

    // Отменить бронирование
    async cancelBooking(eventId) {
        return this.request(`${this.baseURL}/${eventId}`, {
            method: 'POST'
        });
    }
}

// Создаем глобальный экземпляр API
const api = new EventBookerAPI();

// Утилиты для работы с DOM
class DOMUtils {
    // Показать уведомление
    static showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        
        const container = document.querySelector('.app-content') || document.body;
        container.insertBefore(notification, container.firstChild);
        
        // Автоматически скрыть через 5 секунд
        setTimeout(() => {
            notification.remove();
        }, 5000);
    }

    // Показать модальное окно
    static showModal(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.style.display = 'block';
        }
    }

    // Скрыть модальное окно
    static hideModal(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.style.display = 'none';
        }
    }

    // Показать загрузку
    static showLoading(container) {
        const loading = document.createElement('div');
        loading.className = 'loading';
        loading.innerHTML = `
            <div class="spinner"></div>
            <p>Загрузка...</p>
        `;
        container.innerHTML = '';
        container.appendChild(loading);
    }

    // Форматировать дату
    static formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    // Форматировать статус бронирования
    static formatBookingStatus(status) {
        const statusMap = {
            'pending': 'Ожидает подтверждения',
            'confirmed': 'Подтверждено',
            'cancelled': 'Отменено'
        };
        return statusMap[status] || status;
    }

    // Получить CSS класс для статуса
    static getStatusClass(status) {
        return `status-${status}`;
    }
}

// Функции для работы с мероприятиями
class EventManager {
    // Отобразить список мероприятий
    static async renderEvents(container, showBookings = false) {
        try {
            DOMUtils.showLoading(container);
            const events = await api.getEvents();
            
            if (events.length === 0) {
                container.innerHTML = '<p class="loading">Мероприятия не найдены</p>';
                return;
            }

            container.innerHTML = events.map(event => this.createEventCard(event, showBookings)).join('');
            
            // Добавляем обработчики событий
            this.attachEventHandlers(container, showBookings);
            
        } catch (error) {
            container.innerHTML = `<p class="notification error">Ошибка загрузки мероприятий: ${error.message}</p>`;
        }
    }

    // Создать карточку мероприятия
    static createEventCard(event, showBookings = false) {
        const eventDate = DOMUtils.formatDate(event.event_at);
        const bookings = event.bookings || [];
        const confirmedBookings = bookings.filter(b => b.status === 'confirmed').length;
        const pendingBookings = bookings.filter(b => b.status === 'pending').length;
        
        let bookingsHtml = '';
        if (showBookings && bookings.length > 0) {
            bookingsHtml = `
                <div class="bookings-list">
                    <h4>Бронирования:</h4>
                    ${bookings.map(booking => `
                        <div class="booking-card">
                            <div class="booking-info">
                                <div class="booking-id">ID: ${booking.id.substring(0, 8)}...</div>
                                <div class="booking-details">
                                    Мест: ${booking.places_count} | 
                                    Telegram ID: ${booking.telegram_id} | 
                                    Создано: ${DOMUtils.formatDate(booking.created_at)}
                                </div>
                            </div>
                            <div class="booking-status ${DOMUtils.getStatusClass(booking.status)}">
                                ${DOMUtils.formatBookingStatus(booking.status)}
                            </div>
                        </div>
                    `).join('')}
                </div>
            `;
        }

        return `
            <div class="event-card fade-in" data-event-id="${event.id}">
                <div class="event-header">
                    <div>
                        <div class="event-title">${event.title}</div>
                        <div class="event-date">${eventDate}</div>
                    </div>
                </div>
                
                <div class="event-stats">
                    <div class="stat">
                        <div class="stat-value">${event.total_seats}</div>
                        <div class="stat-label">Всего мест</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value">${event.available_seats}</div>
                        <div class="stat-label">Свободно</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value">${confirmedBookings}</div>
                        <div class="stat-label">Подтверждено</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value">${pendingBookings}</div>
                        <div class="stat-label">Ожидает</div>
                    </div>
                </div>

                ${bookingsHtml}

                <div class="event-actions">
                    ${!showBookings ? `
                        <button class="btn btn-primary btn-small" onclick="EventManager.showBookModal('${event.id}')">
                            Забронировать
                        </button>
                        <button class="btn btn-secondary btn-small" onclick="EventManager.showConfirmModal('${event.id}')">
                            Подтвердить бронь
                        </button>
                    ` : `
                        <button class="btn btn-danger btn-small" onclick="EventManager.cancelBooking('${event.id}')">
                            Отменить бронь
                        </button>
                    `}
                </div>
            </div>
        `;
    }

    // Привязать обработчики событий
    static attachEventHandlers(container, showBookings) {
        // Обработчики уже встроены в HTML через onclick
    }

    // Показать модальное окно бронирования
    static showBookModal(eventId) {
        document.getElementById('bookEventId').value = eventId;
        DOMUtils.showModal('bookModal');
    }

    // Показать модальное окно подтверждения
    static showConfirmModal(eventId) {
        document.getElementById('confirmEventId').value = eventId;
        DOMUtils.showModal('confirmModal');
    }

    // Забронировать место
    static async bookEvent() {
        const eventId = document.getElementById('bookEventId').value;
        const telegramId = parseInt(document.getElementById('bookTelegramId').value);
        const placesCount = parseInt(document.getElementById('bookPlacesCount').value);

        if (!telegramId || !placesCount) {
            DOMUtils.showNotification('Заполните все поля', 'error');
            return;
        }

        try {
            await api.bookEvent(eventId, {
                telegram_id: telegramId,
                places_count: placesCount
            });
            
            DOMUtils.showNotification('Место успешно забронировано!', 'success');
            DOMUtils.hideModal('bookModal');
            
            // Обновить список мероприятий
            const container = document.querySelector('.events-list');
            if (container) {
                this.renderEvents(container);
            }
            
        } catch (error) {
            DOMUtils.showNotification(`Ошибка бронирования: ${error.message}`, 'error');
        }
    }

    // Подтвердить бронирование
    static async confirmBooking() {
        const eventId = document.getElementById('confirmEventId').value;

        try {
            await api.confirmBooking(eventId);
            
            DOMUtils.showNotification('Бронирование успешно подтверждено!', 'success');
            DOMUtils.hideModal('confirmModal');
            
            // Обновить список мероприятий
            const container = document.querySelector('.events-list');
            if (container) {
                this.renderEvents(container);
            }
            
        } catch (error) {
            DOMUtils.showNotification(`Ошибка подтверждения: ${error.message}`, 'error');
        }
    }

    // Отменить бронирование
    static async cancelBooking(eventId) {
        if (!confirm('Вы уверены, что хотите отменить это бронирование?')) {
            return;
        }

        try {
            await api.cancelBooking(eventId);
            
            DOMUtils.showNotification('Бронирование успешно отменено!', 'success');
            
            // Обновить список мероприятий
            const container = document.querySelector('.events-list');
            if (container) {
                this.renderEvents(container, true); // Показать с бронированиями для админа
            }
            
        } catch (error) {
            DOMUtils.showNotification(`Ошибка отмены: ${error.message}`, 'error');
        }
    }
}

// Функции для работы с формами
class FormManager {
    // Создать мероприятие
    static async createEvent() {
        const title = document.getElementById('eventTitle').value;
        const eventAtLocal = document.getElementById('eventAt').value;
        const totalSeats = parseInt(document.getElementById('totalSeats').value);

        if (!title || !eventAtLocal || !totalSeats) {
            DOMUtils.showNotification('Заполните все поля', 'error');
            return;
        }

        // Преобразуем datetime-local -> RFC3339
        // datetime-local возвращает дату в локальном времени без зоны
        // Нужно создать Date объект и преобразовать в RFC3339
        const localDate = new Date(eventAtLocal);
        
        // Проверяем, что дата валидна
        if (isNaN(localDate.getTime())) {
            DOMUtils.showNotification('Неверный формат даты', 'error');
            return;
        }

        // Преобразуем в RFC3339 формат
        const eventAt = localDate.toISOString();

        try {
            await api.createEvent({
                title: title,
                event_at: eventAt,
                total_seats: totalSeats
            });
            
            DOMUtils.showNotification('Мероприятие успешно создано!', 'success');
            
            // Очистить форму
            document.getElementById('createEventForm').reset();
            
            // Обновить список мероприятий
            const container = document.querySelector('.events-list');
            if (container) {
                EventManager.renderEvents(container, true); // Показать с бронированиями для админа
            }
            
        } catch (error) {
            DOMUtils.showNotification(`Ошибка создания мероприятия: ${error.message}`, 'error');
        }
    }
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    // Закрытие модальных окон по клику на крестик
    document.querySelectorAll('.close').forEach(closeBtn => {
        closeBtn.addEventListener('click', function() {
            const modal = this.closest('.modal');
            if (modal) {
                modal.style.display = 'none';
            }
        });
    });

    // Закрытие модальных окон по клику вне их
    window.addEventListener('click', function(event) {
        if (event.target.classList.contains('modal')) {
            event.target.style.display = 'none';
        }
    });
});
