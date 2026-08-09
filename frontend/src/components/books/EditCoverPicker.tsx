import { useState, useCallback, useEffect } from 'react';
import { ImagePlus, Upload, X, Trash2, RotateCcw } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

interface EditCoverPickerProps {
  existingCoverUrl?: string | null;
  newFile: File | null;
  shouldDelete: boolean;
  onFileChange: (file: File | null) => void;
  onDeleteChange: (shouldDelete: boolean) => void;
  disabled?: boolean;
}

const ACCEPTED_TYPES = ['image/jpeg', 'image/png', 'image/webp'];
const MAX_SIZE = 5 * 1024 * 1024; // 5MB

export function EditCoverPicker({
  existingCoverUrl,
  newFile,
  shouldDelete,
  onFileChange,
  onDeleteChange,
  disabled,
}: EditCoverPickerProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Generate preview URL when newFile changes
  useEffect(() => {
    if (newFile) {
      const url = URL.createObjectURL(newFile);
      setPreview(url);
      return () => URL.revokeObjectURL(url);
    } else {
      setPreview(null);
    }
  }, [newFile]);

  const validateAndSetFile = useCallback((selectedFile: File) => {
    setError(null);

    if (!ACCEPTED_TYPES.includes(selectedFile.type)) {
      setError('Только JPEG, PNG или WebP');
      return;
    }
    if (selectedFile.size > MAX_SIZE) {
      setError('Файл слишком большой (макс. 5MB)');
      return;
    }

    onFileChange(selectedFile);
    onDeleteChange(false); // If selecting new file, cancel any pending delete
  }, [onFileChange, onDeleteChange]);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    if (disabled) return;

    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) {
      validateAndSetFile(droppedFile);
    }
  }, [disabled, validateAndSetFile]);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    if (!disabled) {
      setIsDragging(true);
    }
  }, [disabled]);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleClick = useCallback(() => {
    if (disabled) return;

    const input = document.createElement('input');
    input.type = 'file';
    input.accept = ACCEPTED_TYPES.join(',');
    input.onchange = (e) => {
      const selectedFile = (e.target as HTMLInputElement).files?.[0];
      if (selectedFile) {
        validateAndSetFile(selectedFile);
      }
    };
    input.click();
  }, [disabled, validateAndSetFile]);

  const handleCancelNewFile = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onFileChange(null);
    setError(null);
  }, [onFileChange]);

  const handleDeleteExisting = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onDeleteChange(true);
    onFileChange(null);
  }, [onDeleteChange, onFileChange]);

  const handleUndoDelete = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onDeleteChange(false);
  }, [onDeleteChange]);

  // State 2 & 5: New file selected - show preview
  if (preview && newFile) {
    return (
      <div className="relative">
        <div className="aspect-[3/4] rounded-lg overflow-hidden bg-muted max-h-48">
          <img
            src={preview}
            alt="Предпросмотр новой обложки"
            className="w-full h-full object-cover"
          />
        </div>
        {!disabled && (
          <Button
            type="button"
            variant="destructive"
            size="icon"
            className="absolute top-1 right-1 h-6 w-6"
            onClick={handleCancelNewFile}
            title="Отменить выбор"
          >
            <X className="h-3 w-3" />
          </Button>
        )}
        <p className="text-xs text-muted-foreground mt-1 text-center truncate">
          {newFile.name}
        </p>
        <p className="text-xs text-primary text-center">Будет загружена</p>
      </div>
    );
  }

  // State 3: Has cover, marked for deletion
  if (shouldDelete && existingCoverUrl) {
    return (
      <div
        className={cn(
          "aspect-[3/4] max-h-48 border-2 border-dashed rounded-lg flex flex-col items-center justify-center transition-colors",
          "border-destructive/50 bg-destructive/5"
        )}
      >
        <div className="flex flex-col items-center gap-2 p-2 text-center">
          <Trash2 className="h-6 w-6 text-destructive" />
          <p className="text-xs text-destructive font-medium">Будет удалена</p>
          {!disabled && (
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="h-7 text-xs"
              onClick={handleUndoDelete}
            >
              <RotateCcw className="h-3 w-3 mr-1" />
              Отменить
            </Button>
          )}
        </div>
      </div>
    );
  }

  // State 1: Has cover, no changes - show existing cover
  if (existingCoverUrl && !shouldDelete) {
    return (
      <div className="relative">
        <div className="aspect-[3/4] rounded-lg overflow-hidden bg-muted max-h-48">
          <img
            src={existingCoverUrl}
            alt="Текущая обложка"
            className="w-full h-full object-cover"
          />
        </div>
        {!disabled && (
          <div className="absolute inset-0 bg-black/60 opacity-0 hover:opacity-100 transition-opacity rounded-lg flex flex-col items-center justify-center gap-2">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="h-7 text-xs"
              onClick={handleClick}
            >
              Заменить
            </Button>
            <Button
              type="button"
              variant="destructive"
              size="sm"
              className="h-7 text-xs"
              onClick={handleDeleteExisting}
            >
              <Trash2 className="h-3 w-3 mr-1" />
              Удалить
            </Button>
          </div>
        )}
      </div>
    );
  }

  // State 4: No cover, no new file - show empty picker
  return (
    <div
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onClick={handleClick}
      className={cn(
        "aspect-[3/4] max-h-48 border-2 border-dashed rounded-lg flex flex-col items-center justify-center cursor-pointer transition-colors",
        isDragging
          ? "border-primary bg-primary/5"
          : "border-border hover:border-primary/50 hover:bg-muted/50",
        disabled && "pointer-events-none opacity-50",
        error && "border-destructive"
      )}
    >
      <div className="flex flex-col items-center gap-1 p-2 text-center">
        <div className="p-2 bg-muted rounded-full">
          {isDragging ? (
            <Upload className="h-4 w-4 text-primary" />
          ) : (
            <ImagePlus className="h-4 w-4 text-muted-foreground" />
          )}
        </div>
        <p className="text-xs font-medium">
          {isDragging ? 'Отпустите' : 'Добавить обложку'}
        </p>
        <p className="text-[10px] text-muted-foreground">
          JPEG, PNG, WebP
        </p>
        {error && (
          <p className="text-[10px] text-destructive">{error}</p>
        )}
      </div>
    </div>
  );
}

