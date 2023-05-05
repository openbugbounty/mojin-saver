function (doc) {
    if (doc.tags) {
        for (tag in doc.tags) {
            if (Array.isArray(doc.tags[tag])) {
                for (val in doc.tags[tag]) {
                    emit(tag, [doc._id, doc.tags[tag][val], doc.program ? doc.program : doc._id], 1)
                }
            } else {
                emit(tag, [doc._id, doc.tags[tag], doc.program ? doc.program : doc._id], 1)
            }
        }
    }
}