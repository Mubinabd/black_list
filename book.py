class Book:
    def __init__(self,author,year,name,isbn) -> None:
        self.author = author
        self.year = year
        self.name = name
        self.isbn = isbn
    def getAuthor (self):
        return self.author
    def getYear(self):
        return self.year
    def getName(self):
        return self.name
    def getIsbn(self):
        return self.isbn
    def toString(self):
        return f'''Author: {self.author},
Name: {self.name},
year: {self.year},
isbn: {self.isbn}'''
book =Book('Tohir Malik',1998,'Shaytanat',52)
print(book.toString())