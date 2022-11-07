
export function setBookmarkData (response, bookmarks) {
  if (response.data.result.length > 0) {
    bookmarks.value = []
    response.data.result.forEach(element => {
      bookmarks.value.push({
        id: element._id,
        folderId: element._source.folderId,
        title: element._source.title,
        date: element._source.date,
        text: element._source.text,
        url: element._source.url,
        favicon: element._source.favicon
      })
    })
  } else {
    bookmarks.value = []
  }
}
