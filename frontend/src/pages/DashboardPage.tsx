import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { TrendingUp, MessageSquare, BookOpen, User, Star } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { StarRating } from "@/components/reviews/StarRating";
import { useDashboard } from "@/api/dashboard";
import { useAuthStore } from "@/stores/auth";
import { formatDistanceToNow } from "@/lib/date";

export function DashboardPage() {
  const { data, isLoading } = useDashboard();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-10 w-48" />
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {Array.from({ length: 6 }).map((_, i) => (
            <Skeleton key={i} className="h-48" />
          ))}
        </div>
      </div>
    );
  }

  const dashboard = data?.data;
  const cacheStatus = data?.cacheStatus;

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-8"
    >
      <div className="flex items-center justify-between">
        <h1 className="font-display text-3xl font-bold">Статистика</h1>
        {cacheStatus && (
          <Badge variant={cacheStatus === "HIT" ? "success" : "outline"}>
            Cache: {cacheStatus}
          </Badge>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Popular Books */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5 text-primary" />
              Популярные книги
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {dashboard?.popular_books.map((book, index) => (
              <Link
                key={book.id}
                to={`/books/${book.id}`}
                className="flex items-center gap-3 p-2 rounded-lg hover:bg-muted transition-colors"
              >
                <span className="text-2xl font-bold text-muted-foreground w-8">
                  {index + 1}
                </span>
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{book.title}</p>
                  <p className="text-sm text-muted-foreground truncate">
                    {book.author}
                  </p>
                </div>
                <div className="flex items-center gap-1 text-sm">
                  <Star className="h-4 w-4 fill-primary text-primary" />
                  <span>{book.average_rating?.toFixed(1) || "-"}</span>
                </div>
              </Link>
            ))}

            {(!dashboard?.popular_books ||
              dashboard.popular_books.length === 0) && (
              <p className="text-center text-muted-foreground py-4">
                Нет данных
              </p>
            )}
          </CardContent>
        </Card>

        {/* Recent Reviews */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <MessageSquare className="h-5 w-5 text-primary" />
              Последние рецензии
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {dashboard?.recent_reviews.map((review) => (
              <Link
                key={review.id}
                to={`/books/${review.book_id}`}
                className="block p-3 rounded-lg hover:bg-muted transition-colors"
              >
                <div className="flex items-center justify-between mb-1">
                  <span className="font-medium text-sm truncate flex-1">
                    {review.book_title}
                  </span>
                  <StarRating value={review.rating} readonly size="sm" />
                </div>
                {review.title && (
                  <p className="text-sm line-clamp-1">{review.title}</p>
                )}
                <div className="flex items-center justify-between mt-2 text-xs text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <User className="h-3 w-3" />
                    {review.user.username}
                  </span>
                  <span>{formatDistanceToNow(review.created_at)}</span>
                </div>
              </Link>
            ))}

            {(!dashboard?.recent_reviews ||
              dashboard.recent_reviews.length === 0) && (
              <p className="text-center text-muted-foreground py-4">
                Нет рецензий
              </p>
            )}
          </CardContent>
        </Card>
      </div>

      {/* User Stats (if authenticated) */}
      {isAuthenticated && dashboard?.user_stats && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <User className="h-5 w-5 text-primary" />
              Ваша статистика
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-4">
              <div className="text-center p-4 bg-muted rounded-lg">
                <BookOpen className="h-8 w-8 mx-auto text-primary mb-2" />
                <p className="text-3xl font-bold">
                  {dashboard.user_stats.books_added}
                </p>
                <p className="text-sm text-muted-foreground">Книг добавлено</p>
              </div>
              <div className="text-center p-4 bg-muted rounded-lg">
                <MessageSquare className="h-8 w-8 mx-auto text-primary mb-2" />
                <p className="text-3xl font-bold">
                  {dashboard.user_stats.reviews_written}
                </p>
                <p className="text-sm text-muted-foreground">
                  Рецензий написано
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </motion.div>
  );
}
