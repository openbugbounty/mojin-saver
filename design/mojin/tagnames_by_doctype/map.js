function (doc) {
    for (prop in doc.tags)
        emit(doc.type, prop);
}