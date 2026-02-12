import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String id;
  final String username;
  final String name;
  final String surname;
  final int role;

  const User({
    required this.id,
    required this.username,
    required this.name,
    required this.surname,
    required this.role,
  });

  bool get isAdmin => role == 1;

  factory User.fromJson(Map<String, dynamic> json) => User(
    id: json['id'] as String,
    username: json['username'] as String,
    name: json['name'] as String,
    surname: (json['surname'] as String?) ?? '',
    role: (json['role'] as int?) ?? 0,
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'username': username,
    'name': name,
    'surname': surname,
    'role': role,
  };

  @override
  List<Object?> get props => [id, username, name, surname, role];
}
