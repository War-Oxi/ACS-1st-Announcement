const API_ENDPOINT = 'http://211.183.3.100:8989/posts';
let currentPage = 1; // 현재 페이지 전역 변수로 선언

function fetchPosts(page = 1) {
  currentPage = page; // 현재 페이지 업데이트
  fetch(`${API_ENDPOINT}?page=${page}`)
    .then(response => response.json())
    .then(data => {
      const posts = data.posts;
      const totalCount = data.total_count; // 총 게시글 수를 받아옵니다
      const totalPages = Math.ceil(totalCount / 15); // 페이지 수를 계산합니다

      const postsTable = document.getElementById('postsTable');
      // Clear table rows except the header
      postsTable.innerHTML = postsTable.rows[0].innerHTML;
      // Populate the posts in the table
      posts.forEach(post => {
        const row = postsTable.insertRow();
        row.innerHTML = `
          <td>${post.post_id}</td>
          <td>${post.post_title}</td>
          <td>${new Date(post.post_create_at).toLocaleString()}</td>
        `;
      });

      setupPagination(totalPages, currentPage); // 페이지네이션 설정
    })
    .catch(error => console.error('Error fetching posts:', error));
}

function setupPagination(totalPages, currentPage = 1) {
  const paginationDiv = document.getElementById('pagination');
  paginationDiv.innerHTML = ''; // Clear the pagination

  // Previous link
  if (currentPage > 1) {
    const prevLink = document.createElement('a');
    prevLink.innerText = '이전';
    prevLink.href = '#';
    prevLink.onclick = () => fetchPosts(currentPage - 1);
    paginationDiv.appendChild(prevLink);
  }

  // Page number links
  for (let i = 1; i <= totalPages; i++) {
    const pageLink = document.createElement('a');
    pageLink.innerText = i;
    pageLink.href = '#';
    pageLink.onclick = () => fetchPosts(i);
    if (i === currentPage) {
      pageLink.className = 'active';
    }
    paginationDiv.appendChild(pageLink);
  }

  // Next link
  if (currentPage < totalPages) {
    const nextLink = document.createElement('a');
    nextLink.innerText = '다음';
    nextLink.href = '#';
    nextLink.onclick = () => fetchPosts(currentPage + 1);
    paginationDiv.appendChild(nextLink);
  }
}

// Initial fetch for the first page
fetchPosts(currentPage);
// Setup dummy pagination (you need to adjust totalPages based on your actual API response)
// setupPagination(10);
