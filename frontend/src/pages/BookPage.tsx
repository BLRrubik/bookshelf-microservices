import { useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { ArrowLeft, Edit, Trash2 } from 'lucide-react';
import { motion } from 'framer-motion';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { StarRating } from '@/components/reviews/StarRating';
import { ReviewCard } from '@/components/reviews/ReviewCard';
import { ReviewForm } from '@/components/reviews/ReviewForm';
import { BookCover } from '@/components/books/BookCover';
import { CoverUpload } from '@/components/books/CoverUpload';
import { CoverStatus } from '@/components/books/CoverStatus';
import { EditCoverPicker } from '@/components/books/EditCoverPicker';
import { useBook, useDeleteBook, useUpdateBook } from '@/api/books';
import { useBookReviews } from '@/api/reviews';
import { useDeleteCover } from '@/api/covers';
import { useAuthStore } from '@/stores/auth';
import { formatDate } from '@/lib/date';
import { gatewayApi } from '@/api/client';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

export function BookPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, user } = useAuthStore();
  
  const { data: book, isLoading: bookLoading } = useBook(id!);
  const { data: reviews, isLoading: reviewsLoading } = useBookReviews(id!);
  const deleteBook = useDeleteBook();
  const updateBook = useUpdateBook(id!);
  const deleteCover = useDeleteCover(id!);

  const [editOpen, setEditOpen] = useState(false);
  const [editTitle, setEditTitle] = useState('');
  const [editAuthor, setEditAuthor] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [editIsbn, setEditIsbn] = useState('');
  const [editYear, setEditYear] = useState('');
  const [editCoverFile, setEditCoverFile] = useState<File | null>(null);
  const [editDeleteCover, setEditDeleteCover] = useState(false);
  const [isEditingCover, setIsEditingCover] = useState(false);

  const isOwner = user && book?.created_by === user.id;

  const handleOpenEdit = () => {
    if (book) {
      setEditTitle(book.title);
      setEditAuthor(book.author);
      setEditDescription(book.description || '');
      setEditIsbn(book.isbn || '');
      setEditYear(book.published_year?.toString() || '');
      setEditCoverFile(null);
      setEditDeleteCover(false);
    }
    setEditOpen(true);
  };

  const handleEdit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      // 1. Update book metadata
      await updateBook.mutateAsync({
        title: editTitle,
        author: editAuthor,
        description: editDescription || undefined,
        isbn: editIsbn || undefined,
        published_year: editYear ? parseInt(editYear) : undefined,
      });

      // 2. Handle cover changes
      if (editDeleteCover && book?.cover_url) {
        setIsEditingCover(true);
        try {
          await deleteCover.mutateAsync();
        } catch {
          toast.warning('Не удалось удалить обложку');
        }
      } else if (editCoverFile) {
        setIsEditingCover(true);
        try {
          const formData = new FormData();
          formData.append('file', editCoverFile);
          await gatewayApi.post(`/api/v1/books/${id}/cover`, formData, {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          });
          toast.success('Обложка загружается...');
        } catch {
          toast.warning('Не удалось загрузить обложку');
        }
      }

      setEditOpen(false);
    } finally {
      setIsEditingCover(false);
    }
  };

  const handleDelete = async () => {
    await deleteBook.mutateAsync(id!);
    navigate('/');
  };

  if (bookLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-32" />
        <div className="grid md:grid-cols-3 gap-8">
          <Skeleton className="aspect-[3/4]" />
          <div className="md:col-span-2 space-y-4">
            <Skeleton className="h-10 w-3/4" />
            <Skeleton className="h-6 w-1/2" />
            <Skeleton className="h-32 w-full" />
          </div>
        </div>
      </div>
    );
  }

  if (!book) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground mb-4">Книга не найдена</p>
        <Button asChild>
          <Link to="/">Вернуться на главную</Link>
        </Button>
      </div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-8"
    >
      <Button variant="ghost" asChild>
        <Link to="/books">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Назад к каталогу
        </Link>
      </Button>

      <div className="grid md:grid-cols-3 gap-8">
        {/* Book Cover */}
        <div className="space-y-4">
          <BookCover
            coverUrl={book.cover_url}
            coverStatus={book.cover_status}
            title={book.title}
            size="lg"
            className="w-full aspect-[3/4]"
          />
          
          <CoverStatus bookId={id!} initialStatus={book.cover_status} />
          
          {/* Cover Upload for Owner */}
          {isOwner && book.cover_status !== 'processing' && (
            <div className="space-y-2">
              {book.cover_url && (
                <Button 
                  variant="outline" 
                  size="sm" 
                  className="w-full"
                  onClick={() => deleteCover.mutate()}
                  disabled={deleteCover.isPending}
                >
                  {deleteCover.isPending ? 'Удаление...' : 'Удалить обложку'}
                </Button>
              )}
              <CoverUpload bookId={id!} />
            </div>
          )}
        </div>

        {/* Book Info */}
        <div className="md:col-span-2 space-y-4">
          <div className="flex items-start justify-between gap-4">
            <div>
              <h1 className="font-display text-3xl font-bold">{book.title}</h1>
              <p className="text-xl text-muted-foreground">{book.author}</p>
            </div>
            
            {isOwner && (
              <div className="flex gap-2">
                <Dialog open={editOpen} onOpenChange={setEditOpen}>
                  <DialogTrigger asChild>
                    <Button variant="outline" size="icon" onClick={handleOpenEdit}>
                      <Edit className="h-4 w-4" />
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="max-w-2xl">
                    <DialogHeader>
                      <DialogTitle>Редактировать книгу</DialogTitle>
                    </DialogHeader>
                    <form onSubmit={handleEdit} className="space-y-4">
                      <div className="grid md:grid-cols-[140px_1fr] gap-4">
                        {/* Cover Picker */}
                        <div className="space-y-2">
                          <Label>Обложка</Label>
                          <EditCoverPicker
                            existingCoverUrl={book.cover_url}
                            newFile={editCoverFile}
                            shouldDelete={editDeleteCover}
                            onFileChange={setEditCoverFile}
                            onDeleteChange={setEditDeleteCover}
                            disabled={updateBook.isPending || isEditingCover}
                          />
                        </div>

                        {/* Form Fields */}
                        <div className="space-y-4">
                          <div className="space-y-2">
                            <Label htmlFor="edit-title">Название</Label>
                            <Input
                              id="edit-title"
                              value={editTitle}
                              onChange={(e) => setEditTitle(e.target.value)}
                              required
                            />
                          </div>
                          <div className="space-y-2">
                            <Label htmlFor="edit-author">Автор</Label>
                            <Input
                              id="edit-author"
                              value={editAuthor}
                              onChange={(e) => setEditAuthor(e.target.value)}
                              required
                            />
                          </div>
                          <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                              <Label htmlFor="edit-isbn">ISBN</Label>
                              <Input
                                id="edit-isbn"
                                value={editIsbn}
                                onChange={(e) => setEditIsbn(e.target.value)}
                              />
                            </div>
                            <div className="space-y-2">
                              <Label htmlFor="edit-year">Год издания</Label>
                              <Input
                                id="edit-year"
                                type="number"
                                value={editYear}
                                onChange={(e) => setEditYear(e.target.value)}
                              />
                            </div>
                          </div>
                        </div>
                      </div>
                      
                      <div className="space-y-2">
                        <Label htmlFor="edit-description">Описание</Label>
                        <Textarea
                          id="edit-description"
                          value={editDescription}
                          onChange={(e) => setEditDescription(e.target.value)}
                          rows={3}
                        />
                      </div>

                      <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => setEditOpen(false)}>
                          Отмена
                        </Button>
                        <Button type="submit" disabled={updateBook.isPending || isEditingCover}>
                          {updateBook.isPending || isEditingCover ? 'Сохранение...' : 'Сохранить'}
                        </Button>
                      </DialogFooter>
                    </form>
                  </DialogContent>
                </Dialog>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button variant="outline" size="icon">
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Удалить книгу?</DialogTitle>
                      <DialogDescription>
                        Это действие нельзя отменить. Книга и все рецензии будут удалены.
                      </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                      <Button variant="outline">Отмена</Button>
                      <Button 
                        variant="destructive" 
                        onClick={handleDelete}
                        disabled={deleteBook.isPending}
                      >
                        {deleteBook.isPending ? 'Удаление...' : 'Удалить'}
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            )}
          </div>

          <div className="flex items-center gap-4">
            <StarRating value={book.average_rating || 0} readonly size="lg" />
            <span className="text-muted-foreground">
              {book.reviews_count} {book.reviews_count === 1 ? 'рецензия' : 'рецензий'}
            </span>
          </div>

          <div className="flex flex-wrap gap-2">
            {book.isbn && <Badge variant="outline">ISBN: {book.isbn}</Badge>}
            {book.published_year && (
              <Badge variant="outline">{book.published_year} г.</Badge>
            )}
          </div>

          {book.description && (
            <p className="text-muted-foreground whitespace-pre-wrap">
              {book.description}
            </p>
          )}

          <p className="text-sm text-muted-foreground">
            Добавлено: {formatDate(book.created_at)}
          </p>
        </div>
      </div>

      <Separator />

      {/* Reviews Section */}
      <div className="space-y-6">
        <h2 className="font-display text-2xl font-semibold">Рецензии</h2>

        {isAuthenticated && (
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Написать рецензию</CardTitle>
            </CardHeader>
            <CardContent>
              <ReviewForm bookId={id!} />
            </CardContent>
          </Card>
        )}

        {reviewsLoading ? (
          <div className="space-y-4">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-32 w-full" />
            ))}
          </div>
        ) : reviews?.data.length === 0 ? (
          <p className="text-center text-muted-foreground py-8">
            Пока нет рецензий. Будьте первым!
          </p>
        ) : (
          <div className="space-y-4">
            {reviews?.data.map((review) => (
              <ReviewCard
                key={review.id}
                review={review}
                bookId={id!}
                currentUserId={user?.id}
              />
            ))}
          </div>
        )}
      </div>
    </motion.div>
  );
}
