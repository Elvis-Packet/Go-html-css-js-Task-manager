async function loadTasks() {
  const res = await fetch('/tasks');
  const tasks = await res.json();
  const list = document.getElementById('task-list');
  list.innerHTML = '';

  tasks.forEach(task => {
    const li = document.createElement('li');
    li.innerHTML = `
      <span>${task.title}</span>
      <button onclick="deleteTask(${task.id})">❌</button>
    `;
    list.appendChild(li);
  });
}

async function addTask(title) {
  await fetch('/tasks', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title })
  });
  loadTasks();
}

async function deleteTask(id) {
  await fetch('/tasks/' + id, { method: 'DELETE' });
  loadTasks();
}

document.getElementById('task-form').addEventListener('submit', e => {
  e.preventDefault();
  const input = document.getElementById('task-input');
  addTask(input.value);
  input.value = '';
});

loadTasks();
