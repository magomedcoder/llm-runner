import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:gen/domain/entities/session.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_bloc.dart';
import 'package:gen/presentation/screens/chat/bloc/chat_event.dart';

Future<void> showSessionTitleEditDialog(
  BuildContext context,
  ChatSession session,
) async {
  final chatBloc = context.read<ChatBloc>();
  final controller = TextEditingController(text: session.title);

  await showDialog<void>(
    context: context,
    builder: (dialogContext) => AlertDialog(
      title: const Text('Редактировать название'),
      content: TextField(
        controller: controller,
        decoration: const InputDecoration(
          hintText: 'Введите новое название',
          border: OutlineInputBorder(),
        ),
        autofocus: true,
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(dialogContext),
          child: const Text('Отмена'),
        ),
        ElevatedButton(
          onPressed: () {
            final title = controller.text.trim();
            if (title.isNotEmpty && title != session.title) {
              chatBloc.add(ChatUpdateSessionTitle(session.id, title));
            }
            Navigator.pop(dialogContext);
          },
          child: const Text('Сохранить'),
        ),
      ],
    ),
  );
}
